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

func (pr *PostRoutes) declarePostRoutes() *chi.Mux {
	postRoute := chi.NewRouter()

	postRoute.Post("/create", middlewares.AuthMiddleware(pr.postController.CreatePost))
	postRoute.Patch("/vote", middlewares.AuthMiddleware(pr.postController.VotePost))
	postRoute.Get("/posts", pr.postController.GetPaginatedPosts)
	postRoute.Get("/{id}", pr.postController.GetPostData)

	return postRoute
}
