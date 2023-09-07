package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/internal/db_types"
	"github.com/adhupraba/breadit-server/lib"
	"github.com/adhupraba/breadit-server/utils"
)

type SubscriptionController struct{}

type subscribeToSubredditBody struct {
	SubredditId int32 `json:"subredditId" validate:"required,gt=0"`
}

type unsubscribeFromSubredditBody struct {
	SubredditId int32 `json:"subredditId" validate:"required,gt=0"`
}

func (sc *SubscriptionController) GetSubredditSubscription(w http.ResponseWriter, r *http.Request, user database.User) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid subreddit id")
		return
	}

	subscription, err := lib.DB.FindUserSubscription(r.Context(), database.FindUserSubscriptionParams{
		UserID:      user.ID,
		SubredditID: int32(id),
	})

	if err != nil && strings.Contains(err.Error(), "no rows") {
		utils.RespondWithJson(w, http.StatusOK, nil)
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, subscription)
}

func (sc *SubscriptionController) SubscribeToSubreddit(w http.ResponseWriter, r *http.Request, user database.User) {
	var body subscribeToSubredditBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	subscription, err := lib.DB.FindUserSubscription(r.Context(), database.FindUserSubscriptionParams{
		UserID:      user.ID,
		SubredditID: body.SubredditId,
	})

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error validating existing subscription status")
		return
	}

	if subscription.ID != 0 {
		utils.RespondWithError(w, http.StatusConflict, "You are already subscribed to this subreddit")
		return
	}

	subscription, err = lib.DB.CreateSubscription(r.Context(), database.CreateSubscriptionParams{
		UserID:      user.ID,
		SubredditID: body.SubredditId,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error while creating subscription")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, true)
}

func (sc *SubscriptionController) UnsubscribeFromSubreddit(w http.ResponseWriter, r *http.Request, user database.User) {
	var body unsubscribeFromSubredditBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	subscription, err := lib.DB.FindUserSubscription(r.Context(), database.FindUserSubscriptionParams{
		UserID:      user.ID,
		SubredditID: body.SubredditId,
	})

	if err != nil && strings.Contains(err.Error(), "no rows") {
		utils.RespondWithError(w, http.StatusConflict, "You are not subscribed to this subreddit inorder to unsubscribe")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error while checking subscription status")
		return
	}

	subreddit, err := lib.DB.FindSubredditOfCreator(r.Context(), database.FindSubredditOfCreatorParams{
		ID:        subscription.SubredditID,
		CreatorID: db_types.NullInt32{Int32: user.ID, Valid: true},
	})

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		fmt.Println("err =>", err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Error while unsubscribing")
		return
	}

	if subreddit.ID != 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "You can't unsubscribe from your own subreddit")
		return
	}

	err = lib.DB.RemoveSubscriptionUsingId(r.Context(), subscription.ID)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error while removing subscription")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, true)
}
