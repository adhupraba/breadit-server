package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/controllers"
)

type UtilsRoutes struct {
	utilsController controllers.UtilsController
}

func GetUtilsRoutes() *chi.Mux {
	routes := UtilsRoutes{}
	return routes.declareUtilsRoutes()
}

func (ur *UtilsRoutes) declareUtilsRoutes() *chi.Mux {
	utilsRoute := chi.NewRouter()

	utilsRoute.Get("/link", ur.utilsController.GetUrlMetadata)

	return utilsRoute
}
