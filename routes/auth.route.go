package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/controllers"
)

type AuthRoutes struct {
	authController controllers.AuthController
}

func GetAuthRoutes() *chi.Mux {
	routes := AuthRoutes{}
	return routes.declareAuthRoutes()
}

func (ar *AuthRoutes) declareAuthRoutes() *chi.Mux {
	authRoute := chi.NewRouter()

	authRoute.Post("/sign-up", ar.authController.Signup)
	authRoute.Post("/sign-in", ar.authController.Signin)

	return authRoute
}
