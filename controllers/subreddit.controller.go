package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/internal/types"
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
		CreatorID: types.NullInt32{Int32: user.ID, Valid: true},
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

func (sc *SubredditController) GetSubredditDataWithErrors(w http.ResponseWriter, r *http.Request) {
	subredditName := chi.URLParam(r, "name")

	subreddit, err := lib.DB.FindSubredditByName(r.Context(), subredditName)

	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Subreddit not found.")
		return
	}

	posts, err := lib.DB.FindPostsOfASubreddit(r.Context(), database.FindPostsOfASubredditParams{
		SubredditID: subreddit.ID,
		Offset:      0,
		Limit:       10,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching posts from the subreddit.")
		return
	}

	type postWithData struct {
		database.Post
		Author   database.User      `json:"author"`
		Votes    []database.Vote    `json:"votes,omitempty"`
		Comments []database.Comment `json:"comments,omitempty"`
	}

	var postsData []postWithData

	for _, data := range posts {
		var votes []database.Vote
		var comments []database.Comment

		fmt.Printf("votes %v\n", votes)

		json.Unmarshal(data.Votes, &votes)
		json.Unmarshal(data.Comments, &comments)

		postsData = append(postsData, postWithData{
			database.Post{
				ID:          data.ID,
				Title:       data.Title,
				Content:     data.Content,
				SubredditID: data.SubredditID,
				AuthorID:    data.AuthorID,
				CreatedAt:   data.CreatedAt,
				UpdatedAt:   data.UpdatedAt,
			},
			data.User,
			votes,
			comments,
		})
	}

	type response struct {
		database.Subreddit
		Posts []postWithData `json:"posts"`
	}

	utils.RespondWithJson(w, http.StatusOK, response{
		subreddit,
		postsData,
	})

	// type response struct {
	// 	database.Subreddit
	// 	Posts []database.FindPostsOfASubredditRow `json:"posts"`
	// }

	// utils.RespondWithJson(w, http.StatusOK, response{
	// 	subreddit,
	// 	posts,
	// })
}
