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
		// ! by adding `127.0.0.1` in the Addr field, the popup to accept incoming connections will not appear anymore
		// ! note that the client also must have the api url as http://127.0.0.1 instead of http://localhost
		// ! if http://localhost was used, then the client will try to connect to ::1 on whatever port used (::1 is the ipv6 of localhost)
		// Addr: "127.0.0.1:" + lib.EnvConfig.Port,
		Addr: ":" + lib.EnvConfig.Port,
	}

	log.Printf("Server listening on port %s", lib.EnvConfig.Port)

	apiRouter := chi.NewRouter()

	apiRouter.Mount("/health", routes.GetHealthRoutes())
	apiRouter.Mount("/auth", routes.GetAuthRoutes())
	apiRouter.Mount("/subreddit", routes.GetSubredditRoutes())
	apiRouter.Mount("/subscription", routes.GetSubscriptionRoutes())
	apiRouter.Mount("/post", routes.GetPostRoutes())
	apiRouter.Mount("/comment", routes.GetCommentRoutes())
	apiRouter.Mount("/utils", routes.GetUtilsRoutes())
	apiRouter.Mount("/search", routes.GetSearchRoutes())

	router.Mount("/api", apiRouter)

	err := serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
