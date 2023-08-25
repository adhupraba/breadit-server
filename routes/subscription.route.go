package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/controllers"
	"github.com/adhupraba/breadit-server/middlewares"
)

type SubscriptionRoutes struct {
	subscriptionController controllers.SubscriptionController
}

func GetSubscriptionRoutes() *chi.Mux {
	routes := SubscriptionRoutes{}
	return routes.declareSubscriptionRoutes()
}

func (sr *SubscriptionRoutes) declareSubscriptionRoutes() *chi.Mux {
	subscriptionRoute := chi.NewRouter()

	subscriptionRoute.Get("/subreddit/{id}", middlewares.AuthMiddleware(sr.subscriptionController.GetSubredditSubscription))
	subscriptionRoute.Post("/subscribe", middlewares.AuthMiddleware(sr.subscriptionController.SubscribeToSubreddit))
	subscriptionRoute.Post("/unsubscribe", middlewares.AuthMiddleware(sr.subscriptionController.UnsubscribeFromSubreddit))

	return subscriptionRoute
}
