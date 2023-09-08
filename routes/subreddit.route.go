package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/controllers"
	"github.com/adhupraba/breadit-server/middlewares"
)

type SubredditRoutes struct {
	subRedditController controllers.SubredditController
}

func GetSubredditRoutes() *chi.Mux {
	routes := SubredditRoutes{}
	return routes.declareSubredditRoutes()
}

func (sr *SubredditRoutes) declareSubredditRoutes() *chi.Mux {
	subredditRoute := chi.NewRouter()

	subredditRoute.Post("/", middlewares.AuthMiddleware(sr.subRedditController.CreateSubreddit))
	subredditRoute.Get("/{name}", sr.subRedditController.GetSubredditDataWithPosts)
	subredditRoute.Get("/list", sr.subRedditController.GetPaginatedSubredditList)

	return subredditRoute
}
