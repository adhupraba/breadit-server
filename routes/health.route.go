package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/controllers"
)

type HealthRoutes struct {
	healthController controllers.HealthController
}

func GetHealthRoutes() *chi.Mux {
	routes := HealthRoutes{}
	return routes.declareHealthRoutes()
}

func (hr *HealthRoutes) declareHealthRoutes() *chi.Mux {
	healthRoute := chi.NewRouter()

	healthRoute.Get("/heartbeat", hr.healthController.Heartbeat)

	return healthRoute
}
