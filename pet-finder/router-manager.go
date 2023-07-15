package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/younesabouali/formal-challenges/pet-finder/Controllers"
	"github.com/younesabouali/formal-challenges/pet-finder/Middlewares"
	"github.com/younesabouali/formal-challenges/pet-finder/internal/database"
)

func AppRouter(port string, DB *database.Queries) {

	router := chi.NewRouter()
	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"*"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           300,
			},
		),
	)
	middlewares := Middlewares.Middlewares{DB: DB}
	routeAuth, store := middlewares.Init()
	v1Router := chi.NewRouter()
	v1Router.Mount("/users", Controllers.UserRouter(DB, middlewares))
	v1Router.Mount("/missing_pets", Controllers.MissingPetsRouter(DB, middlewares))
	v1Router.Mount("/events_recorder", Controllers.EventsRecorderRouter(DB, middlewares))
	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	fmt.Println("server runing on PORT : ", port)
	srv.ListenAndServe()
}
