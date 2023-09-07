package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/controllers"
	"github.com/adhupraba/breadit-server/middlewares"
)

type CommentRoutes struct {
	commentController controllers.CommentController
}

func GetCommentRoutes() *chi.Mux {
	routes := CommentRoutes{}
	return routes.declareCommentRoutes()
}

func (cr *CommentRoutes) declareCommentRoutes() *chi.Mux {
	commentRoute := chi.NewRouter()

	commentRoute.Post("/create", middlewares.AuthMiddleware(cr.commentController.CreateComment))
	commentRoute.Get("/post/{postId}", cr.commentController.GetCommentsOfAPost)
	commentRoute.Patch("/vote", middlewares.AuthMiddleware(cr.commentController.VoteComment))

	return commentRoute
}
