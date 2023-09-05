package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/adhupraba/breadit-server/constants"
	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/internal/helpers"
	"github.com/adhupraba/breadit-server/internal/helpers/transformer"
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

	content := types.NullRawMessage{}

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

	post, err := lib.DB.FindPostWithAuthorAndVotes(r.Context(), body.PostId)

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when fetching post details")
		return
	}

	if post.Post.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	transformedPost, err := transformer.TransformPostWithAuthorAndVotes(post)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	existingVote, err := lib.DB.FindUserVoteOfAPost(r.Context(), database.FindUserVoteOfAPostParams{
		UserID: user.ID,
		PostID: body.PostId,
	})

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		utils.RespondWithError(w, http.StatusBadRequest, "Error while voting the post")
		return
	}

	// existing vote is present
	if existingVote.ID != 0 {
		// existing vote is of the same type, delete the vote
		if existingVote.Type == body.VoteType {
			err = lib.DB.RemoveVote(r.Context(), existingVote.ID)

			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error when removing vote")
				return
			}

			utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "Successfully removed vote"})
			return
		}

		// existing vote is of different type, update the vote
		err = lib.DB.UpdateVote(r.Context(), database.UpdateVoteParams{
			ID:   existingVote.ID,
			Type: body.VoteType,
		})

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error when updating vote")
			return
		}

		// recount the votes
		voteCount := 0

		for _, vote := range transformedPost.Votes {
			if vote.Type == database.VoteTypeUP {
				voteCount += 1
			} else if vote.Type == database.VoteTypeDOWN {
				voteCount -= 1
			}
		}

		if voteCount >= constants.CacheAfterUpvotes {
			cachePayload := helpers.RedisCachedPost{
				ID:             transformedPost.ID,
				Title:          transformedPost.Title,
				AuthorUsername: transformedPost.Author.Username,
				Content:        transformedPost.Content,
				CurrentVote:    database.NullVoteType{VoteType: body.VoteType, Valid: true},
				CreatedAt:      transformedPost.CreatedAt,
			}

			key := fmt.Sprintf("post:%d", transformedPost.ID)
			lib.Redis.HSet(r.Context(), key, cachePayload)
		}

		utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "Successfully updated vote on post"})
		return
	}

	// create new vote
	_, err = lib.DB.CreateVote(r.Context(), database.CreateVoteParams{
		PostID: body.PostId,
		UserID: user.ID,
		Type:   body.VoteType,
	})

	// recount the votes
	voteCount := 0

	for _, vote := range transformedPost.Votes {
		if vote.Type == database.VoteTypeUP {
			voteCount += 1
		} else if vote.Type == database.VoteTypeDOWN {
			voteCount -= 1
		}
	}

	if voteCount >= constants.CacheAfterUpvotes {
		cachePayload := helpers.RedisCachedPost{
			ID:             transformedPost.ID,
			Title:          transformedPost.Title,
			AuthorUsername: transformedPost.Author.Username,
			Content:        transformedPost.Content,
			CurrentVote:    database.NullVoteType{VoteType: body.VoteType, Valid: true},
			CreatedAt:      transformedPost.CreatedAt,
		}

		key := fmt.Sprintf("post:%d", transformedPost.ID)
		lib.Redis.HSet(r.Context(), key, cachePayload)
	}

	utils.RespondWithJson(w, http.StatusCreated, types.Json{"message": "Successfully added vote on post"})
}

func (pc *PostController) GetPaginatedPosts(w http.ResponseWriter, r *http.Request) {
	type urlSearchParams struct {
		Limit         int    `validate:"required,number,gt=0"`
		Page          int    `validate:"required,number,gt=0"`
		SubredditName string `validate:"min=3"`
	}

	query := r.URL.Query()

	limit, err := strconv.Atoi(query.Get("limit"))

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Recevied invalid limit value")
		return
	}

	page, err := strconv.Atoi(query.Get("page"))

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Recevied invalid page value")
		return
	}

	searchParams := urlSearchParams{
		Limit:         limit,
		Page:          page,
		SubredditName: query.Get("subredditName"),
	}

	err = lib.Validate.Struct(searchParams)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Recevied invalid ")
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

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error getting posts")
		return
	}

	// ! postsWithData can be `nil` if data is not found for a given pagination.
	// ! eg: only 10 posts are there. and now i am requesting 11-15 posts. it will be `nil`
	postsWithData, err, errCode := helpers.GetPostsOfSubreddit(r.Context(), params)

	if err != nil {
		utils.RespondWithError(w, errCode, err.Error())
		return
	}

	if postsWithData == nil {
		postsWithData = []helpers.PostWithData{}
	}

	utils.RespondWithJson(w, http.StatusOK, postsWithData)
}
