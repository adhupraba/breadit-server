package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	"github.com/adhupraba/breadit-server/lib"
	"github.com/adhupraba/breadit-server/routes"
)

func init() {
	lib.LoadEnv()
	lib.ConnectDb()
	lib.ConnectRedis()
}

func main() {
	router := chi.NewRouter()
	router.Use(cors.AllowAll().Handler)

	serve := http.Server{
		Handler: router,
		Addr:    ":" + lib.EnvConfig.Port,
	}

	log.Printf("Server listening on port %s", lib.EnvConfig.Port)

	apiRouter := chi.NewRouter()

	apiRouter.Mount("/health", routes.GetHealthRoutes())
	apiRouter.Mount("/auth", routes.GetAuthRoutes())

	router.Mount("/api", apiRouter)

	err := serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
