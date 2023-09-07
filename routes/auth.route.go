package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/controllers"
	"github.com/adhupraba/breadit-server/middlewares"
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
	authRoute.Get("/sign-out", ar.authController.LogoutUser)
	authRoute.Get("/refresh", ar.authController.RefreshAccessToken)
	authRoute.Get("/get-me", middlewares.AuthMiddleware(ar.authController.GetUser))
	authRoute.Patch("/username", middlewares.AuthMiddleware(ar.authController.UpdateUsername))

	return authRoute
}
