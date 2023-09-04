package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/controllers"
	"github.com/adhupraba/breadit-server/middlewares"
)

type PostRoutes struct {
	postController controllers.PostController
}

func GetPostRoutes() *chi.Mux {
	routes := PostRoutes{}
	return routes.declarePostRoutes()
}

func (sr *PostRoutes) declarePostRoutes() *chi.Mux {
	postRoute := chi.NewRouter()

	postRoute.Post("/create", middlewares.AuthMiddleware(sr.postController.CreatePost))
	postRoute.Patch("/vote", middlewares.AuthMiddleware(sr.postController.VotePost))

	return postRoute
}
