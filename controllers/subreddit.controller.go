package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/internal/db_types"
	queryhelpers "github.com/adhupraba/breadit-server/internal/helpers/query_helpers"
	"github.com/adhupraba/breadit-server/lib"
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
		CreatorID: db_types.NullInt32{Int32: user.ID, Valid: true},
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

	utils.RespondWithJson(w, http.StatusCreated, subreddit)
}

func (sc *SubredditController) GetSubredditDataWithPosts(w http.ResponseWriter, r *http.Request) {
	subredditName := chi.URLParam(r, "name")

	data, err, errCode := queryhelpers.GetSubredditWithPosts(r.Context(), subredditName)

	if err != nil {
		utils.RespondWithError(w, errCode, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, data)
}

func (sc *SubredditController) GetPaginatedSubredditList(w http.ResponseWriter, r *http.Request) {
	type urlSearchParams struct {
		Limit     int `validate:"required,number,gt=0"`
		Page      int `validate:"required,number,gte=0"`
		CreatedBy string
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
		Limit:     limit,
		Page:      page,
		CreatedBy: query.Get("createdBy"),
	}

	var userId db_types.NullInt32
	cookie, _ := r.Cookie("access_token")

	if cookie != nil && searchParams.CreatedBy == "self" {
		user, _ := utils.GetUserFromToken(w, r, cookie.Value)

		if user.ID != 0 {
			userId = db_types.NullInt32{
				Int32: user.ID,
				Valid: true,
			}
		}
	}

	subreddits, err := lib.DB.SearchSubreddits(r.Context(), database.SearchSubredditsParams{
		Offset: (int32(searchParams.Page) - 1) * int32(searchParams.Limit),
		Limit:  int32(searchParams.Limit),
		Name:   "%",
		UserID: userId,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching subreddits")
		return
	}

	if subreddits == nil {
		subreddits = []database.SearchSubredditsRow{}
	}

	utils.RespondWithJson(w, http.StatusOK, subreddits)
}
