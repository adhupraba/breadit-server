package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	"github.com/adhupraba/breadit-server/lib"
	"github.com/adhupraba/breadit-server/routes"
)

var embedMigrations embed.FS

func init() {
	lib.LoadEnv()
	lib.ConnectDb()
	lib.ConnectRedis()
}

func main() {
	router := chi.NewRouter()
	// router.Use(cors.AllowAll().Handler)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	serve := http.Server{
		Handler: router,
		Addr:    ":" + lib.EnvConfig.Port,
	}

	log.Printf("Server listening on port %s", lib.EnvConfig.Port)

	apiRouter := chi.NewRouter()

	apiRouter.Mount("/health", routes.GetHealthRoutes())
	apiRouter.Mount("/auth", routes.GetAuthRoutes())
	apiRouter.Mount("/subreddit", routes.GetSubredditRoutes())
	apiRouter.Mount("/subscription", routes.GetSubscriptionRoutes())
	apiRouter.Mount("/post", routes.GetPostRoutes())
	apiRouter.Mount("/utils", routes.GetUtilsRoutes())

	router.Mount("/api", apiRouter)

	err := serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
