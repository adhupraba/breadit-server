package middlewares

import (
	"net/http"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/utils"
)

type NextFunc func(http.ResponseWriter, *http.Request, database.User)

func AuthMiddleware(handler NextFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")

		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, "Could not get access token")
			return
		}

		user, err := utils.GetUserFromToken(w, r, cookie.Value)

		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, err.Error())
			return
		}

		handler(w, r, user)
	}
}
