package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/internal/db_types"
	queryhelpers "github.com/adhupraba/breadit-server/internal/helpers/query_helpers"
	"github.com/adhupraba/breadit-server/internal/types"
	"github.com/adhupraba/breadit-server/lib"
	"github.com/adhupraba/breadit-server/utils"
)

type CommentController struct{}

type createCommentBody struct {
	PostId    int32              `json:"postId" validate:"required,gt=0"`
	Text      string             `json:"text" validate:"required,min=1"`
	ReplyToId db_types.NullInt32 `json:"replyToId"`
}

type voteCommentBody struct {
	CommentId int32             `json:"commentId" validate:"required,gt=0"`
	VoteType  database.VoteType `json:"voteType" validate:"required,oneof=UP DOWN"`
}

func (cc *CommentController) CreateComment(w http.ResponseWriter, r *http.Request, user database.User) {
	var body createCommentBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	newComment, err := lib.DB.CreateComment(r.Context(), database.CreateCommentParams{
		Text:      body.Text,
		PostID:    body.PostId,
		AuthorID:  user.ID,
		ReplyToID: body.ReplyToId,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when creating comment")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, newComment)
}

func (cc *CommentController) GetCommentsOfAPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "postId")
	postId, err := strconv.Atoi(id)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid post id")
		return
	}

	comments, err, errCode := queryhelpers.FindCommentsAndRepliesOfAPost(r.Context(), int32(postId))

	if err != nil {
		utils.RespondWithError(w, errCode, "Error when fetching comments")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, comments)
}

func (cc *CommentController) VoteComment(w http.ResponseWriter, r *http.Request, user database.User) {
	var body voteCommentBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	existingVote, err := lib.DB.FindUserVoteOfAComment(r.Context(), database.FindUserVoteOfACommentParams{
		UserID:    user.ID,
		CommentID: body.CommentId,
	})

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		utils.RespondWithError(w, http.StatusBadRequest, "Error when checking for existing vote on comment")
		return
	}

	// existing vote is present
	if existingVote.ID != 0 {
		if existingVote.Type == body.VoteType {
			err = lib.DB.RemoveCommentVote(r.Context(), body.CommentId)

			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error when removing vote")
				return
			}

			utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "Successfully removed vote"})
			return
		}

		// existing vote is of different type, update the vote
		err = lib.DB.UpdateCommentVote(r.Context(), database.UpdateCommentVoteParams{
			ID:   existingVote.ID,
			Type: body.VoteType,
		})

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error when updating vote")
			return
		}

		utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "Successfully updated vote on comment"})
		return
	}

	_, err = lib.DB.CreateCommentVote(r.Context(), database.CreateCommentVoteParams{
		CommentID: body.CommentId,
		UserID:    user.ID,
		Type:      body.VoteType,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating comment")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, types.Json{"message": "Successfully voted on comment"})
}
