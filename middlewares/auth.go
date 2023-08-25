package middlewares

import (
	"fmt"
	"net/http"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/utils"
)

type NextFunc func(http.ResponseWriter, *http.Request, database.User)

func AuthMiddleware(handler NextFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")

		if err != nil {
			fmt.Println(r.URL, "Could not get access token =>", err.Error())
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized access")
			return
		}

		user, err := utils.GetUserFromToken(w, r, cookie.Value)

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		handler(w, r, user)
	}
}
