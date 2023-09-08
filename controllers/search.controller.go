package controllers

import (
	"net/http"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/lib"
	"github.com/adhupraba/breadit-server/utils"
)

type SearchController struct{}

func (cc *SearchController) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid query")
		return
	}

	results, err := lib.DB.SearchSubreddits(r.Context(), database.SearchSubredditsParams{
		Name:   query + "%",
		Offset: 0,
		Limit:  5,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when searching subreddits")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, results)
}
