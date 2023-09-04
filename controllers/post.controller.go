package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/adhupraba/breadit-server/internal/database"
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

func (sc *PostController) CreatePost(w http.ResponseWriter, r *http.Request, user database.User) {
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

func (sc *PostController) VotePost(w http.ResponseWriter, r *http.Request, user database.User) {
	var body votePostBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	post, err := lib.DB.FindPostWithAuthorAndVotes(r.Context(), body.PostId)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when fetching post details")
		return
	}

	transformedPost, err := transformer.TransformPostWithAuthorAndVotes(post)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, transformedPost)

	// votes, err := lib.DB.FindVotesOfAPost(r.Context(), body.PostId)

	// if err != nil {
	// 	utils.RespondWithError(w, http.StatusInternalServerError, "Error when fetching votes of the post")
	// 	return
	// }

	// existingVote, err := lib.DB.FindUserVoteOfAPost(r.Context(), database.FindUserVoteOfAPostParams{
	// 	UserId: user.ID,
	// 	PostId: body.PostId,
	// })

	// if err != nil && !strings.Contains(err.Error(), "no rows") {
	// 	utils.RespondWithError(w, http.StatusBadRequest, "Error while voting the post")
	// 	return
	// }

	// // existing vote present, check and update the vote
	// if existingVote.ID != 0 {
	// 	if existingVote.Type == body.VoteType {
	// 		utils.RespondWithError(w, http.StatusBadRequest, "You cannot vote the same type twice")
	// 		return
	// 	}

	// 	err = lib.DB.UpdateVote(r.Context(), database.UpdateVoteParams{
	// 		ID:   existingVote.ID,
	// 		Type: body.VoteType,
	// 	})

	// 	utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "Successfully updated vote"})
	// 	return
	// }

	// // create new vote
	// _, err = lib.DB.CreateVote(r.Context(), database.CreateVoteParams{
	// 	PostId: body.PostId,
	// 	UserId: user.ID,
	// 	Type:   body.VoteType,
	// })

	// utils.RespondWithJson(w, http.StatusCreated, types.Json{"message": "Successfully voted on post"})
}
