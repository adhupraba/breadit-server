package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/breadit-server/controllers"
)

type SearchRoutes struct {
	searchController controllers.SearchController
}

func GetSearchRoutes() *chi.Mux {
	routes := SearchRoutes{}
	return routes.declareSearchRoutes()
}

func (cr *SearchRoutes) declareSearchRoutes() *chi.Mux {
	searchRoute := chi.NewRouter()

	searchRoute.Get("/", cr.searchController.Search)

	return searchRoute
}
