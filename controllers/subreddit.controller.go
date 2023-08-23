package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/lib"
	"github.com/adhupraba/breadit-server/models"
	"github.com/adhupraba/breadit-server/utils"
)

type SubredditController struct{}

type createSubredditBody struct {
	Name string `json:"name" validate:"required,min=3,max=21"`
}

func (sc *SubredditController) CreateSubreddit(w http.ResponseWriter, r *http.Request, user database.User) {
	var body createSubredditBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	subreddit, err := lib.DB.FindSubredditByName(r.Context(), body.Name)

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		fmt.Println("existing subreddit db error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Unable to validate subreddit name.")
		return
	}

	if subreddit.ID != 0 {
		utils.RespondWithError(w, http.StatusConflict, "Subreddit name already exists.")
		return
	}

	subreddit, err = lib.DB.CreateSubreddit(r.Context(), database.CreateSubredditParams{
		Name:      body.Name,
		CreatorID: sql.NullInt32{Int32: user.ID, Valid: true},
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to create subreddit.")
		return
	}

	_, err = lib.DB.CreateSubscription(r.Context(), database.CreateSubscriptionParams{
		UserID:      user.ID,
		SubredditID: subreddit.ID,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to subscribe to created subreddit.")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, models.DbSubredditToSubreddit(subreddit))
}

func (sc *SubredditController) GetSubredditDataWithErrors(w http.ResponseWriter, r *http.Request) {
	// subredditName := chi.URLParam(r, "name")

	// type response struct {
	// 	Subreddit database.Subreddit

	// }
}
