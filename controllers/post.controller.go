package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/constants"
	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/internal/db_types"
	"github.com/adhupraba/breadit-server/internal/helpers"
	queryhelpers "github.com/adhupraba/breadit-server/internal/helpers/query_helpers"
	"github.com/adhupraba/breadit-server/internal/types"
	"github.com/adhupraba/breadit-server/lib"
	"github.com/adhupraba/breadit-server/utils"
)

type PostController struct{}

type createPostBody struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Content     any    `json:"content"`
	SubredditId int32  `json:"subredditId" validate:"required,gt=0"`
}

type votePostBody struct {
	PostId   int32             `json:"postId" validate:"required,gt=0"`
	VoteType database.VoteType `json:"voteType" validate:"oneof=UP DOWN"`
}

func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request, user database.User) {
	var body createPostBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, err = lib.DB.FindUserSubscription(r.Context(), database.FindUserSubscriptionParams{
		UserID:      user.ID,
		SubredditID: body.SubredditId,
	})

	if err != nil && strings.Contains(err.Error(), "no rows") {
		utils.RespondWithError(w, http.StatusBadRequest, "Subscribe to post")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	content := db_types.NullRawMessage{}

	if body.Content != nil {
		contentJson, err := json.Marshal(body.Content)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Unable to parse the post content")
			return
		}

		content.RawMessage = contentJson
		content.Valid = true
	}

	_, err = lib.DB.CreatePost(r.Context(), database.CreatePostParams{
		Title:       body.Title,
		Content:     content,
		SubredditID: body.SubredditId,
		AuthorID:    user.ID,
	})

	utils.RespondWithJson(w, http.StatusCreated, types.Json{"message": "Successfully created post"})
}

func (pc *PostController) VotePost(w http.ResponseWriter, r *http.Request, user database.User) {
	var body votePostBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	post, err, errCode := queryhelpers.FindPostWithAuthorAndVotes(r.Context(), body.PostId)

	if err != nil {
		utils.RespondWithError(w, errCode, err.Error())
		return
	}

	existingVote, err := lib.DB.FindUserVoteOfAPost(r.Context(), database.FindUserVoteOfAPostParams{
		UserID: user.ID,
		PostID: body.PostId,
	})

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		utils.RespondWithError(w, http.StatusBadRequest, "Error when checking for existing vote on post")
		return
	}

	// existing vote is present
	if existingVote.ID != 0 {
		// existing vote is of the same type, delete the vote
		if existingVote.Type == body.VoteType {
			err = lib.DB.RemovePostVote(r.Context(), existingVote.ID)

			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error when removing vote")
				return
			}

			utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "Successfully removed vote"})
			return
		}

		// existing vote is of different type, update the vote
		err = lib.DB.UpdatePostVote(r.Context(), database.UpdatePostVoteParams{
			ID:   existingVote.ID,
			Type: body.VoteType,
		})

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error when updating vote")
			return
		}

		recountVotesAndSaveInRedis(r.Context(), post)

		utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "Successfully updated vote on post"})
		return
	}

	// create new vote
	_, err = lib.DB.CreatePostVote(r.Context(), database.CreatePostVoteParams{
		PostID: body.PostId,
		UserID: user.ID,
		Type:   body.VoteType,
	})

	recountVotesAndSaveInRedis(r.Context(), post)

	utils.RespondWithJson(w, http.StatusCreated, types.Json{"message": "Successfully voted on post"})
}

func (pc *PostController) GetPaginatedPosts(w http.ResponseWriter, r *http.Request) {
	type urlSearchParams struct {
		Limit         int `validate:"required,number,gt=0"`
		Page          int `validate:"required,number,gte=0"`
		SubredditName string
	}

	query := r.URL.Query()

	limit, err := strconv.Atoi(query.Get("limit"))

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Received invalid limit value")
		return
	}

	page, err := strconv.Atoi(query.Get("page"))

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Received invalid page value")
		return
	}

	searchParams := urlSearchParams{
		Limit:         limit,
		Page:          page,
		SubredditName: query.Get("subredditName"),
	}

	err = lib.Validate.Struct(searchParams)

	if err != nil {
		fmt.Println("search params validation error =>", err)
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Received invalid search params")
		return
	}

	followingSubredditIds := []int32{}
	isAuthenticated := false

	cookie, _ := r.Cookie("access_token")

	if cookie != nil {
		user, _ := utils.GetUserFromToken(w, r, cookie.Value)

		if user.ID != 0 {
			isAuthenticated = true
			subscriptions, err := lib.DB.FindAllSubscriptionsOfUser(r.Context(), user.ID)

			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error getting subscribed subreddits")
				return
			}

			for _, sub := range subscriptions {
				followingSubredditIds = append(followingSubredditIds, sub.ID)
			}
		}
	}

	params := database.FindPostsOfSubredditParams{
		SubredditName:   searchParams.SubredditName,
		IsAuthenticated: isAuthenticated,
		SubredditIds:    followingSubredditIds,
		Offset:          (int32(searchParams.Page) - 1) * int32(searchParams.Limit),
		Limit:           int32(searchParams.Limit),
	}
	// fmt.Printf("GetPaginatedPosts db params => %#v\n", params)

	// ! postsWithData can be `nil` if data is not found for a given pagination.
	// ! eg: only 10 posts are there. and now i am requesting 11-15 posts. it will be `nil`
	postsWithData, err, errCode := queryhelpers.GetPostsOfSubreddit(r.Context(), params)

	if err != nil {
		utils.RespondWithError(w, errCode, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, postsWithData)
}

func (pc *PostController) GetPostData(w http.ResponseWriter, r *http.Request) {
	type response struct {
		CachedPost *helpers.RedisCachedPost      `json:"cachedPost"`
		Post       *types.PostWithAuthorAndVotes `json:"post"`
	}

	id := chi.URLParam(r, "id")
	postId, err := strconv.Atoi(id)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid post id")
		return
	}

	isForceFetch := r.URL.Query().Get("force") == "true"

	if !isForceFetch {
		key := fmt.Sprintf("post:%d", postId)
		cachedPost, err := helpers.GetPostFromRedis(r.Context(), key)

		if err == nil && cachedPost.ID != 0 {
			utils.RespondWithJson(w, http.StatusOK, response{CachedPost: &cachedPost, Post: nil})
			return
		}
	}

	post, err, errCode := queryhelpers.FindPostWithAuthorAndVotes(r.Context(), int32(postId))

	if err != nil {
		utils.RespondWithError(w, errCode, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, response{CachedPost: nil, Post: &post})
}

// ! re-usable helpers
func recountVotesAndSaveInRedis(ctx context.Context, post types.PostWithAuthorAndVotes) {
	// recount the votes
	voteCount := 0

	for _, vote := range post.Votes {
		if vote.Type == database.VoteTypeUP {
			voteCount += 1
		} else if vote.Type == database.VoteTypeDOWN {
			voteCount -= 1
		}
	}

	fmt.Println("voteCount", voteCount)

	if voteCount >= constants.CacheAfterUpvotes {
		cachePayload := helpers.RedisCachedPost{
			ID:             post.ID,
			Title:          post.Title,
			AuthorUsername: post.Author.Username,
			Content:        post.Content,
			CreatedAt:      post.CreatedAt,
		}

		key := fmt.Sprintf("post:%d", post.ID)
		err := helpers.SavePostInRedis(ctx, key, cachePayload)

		if err != nil {
			fmt.Println("redis hset error =>", err)
		}

		// cachedPost, err := helpers.GetPostFromRedis(r.Context(), key)

		// if err != nil {
		// 	fmt.Println("redis hgetall error =>", err)
		// }

		// fmt.Printf("cached post from redis => %#v\n", cachedPost)
	}
}
