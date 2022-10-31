package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/exercises", app.GetExercises)
	mux.Get("/exercises/{date}/{name}", app.GetByNameDate)
	mux.Get("/exercises/name/{name}", app.GetByName)
	mux.Get("/exercises/date/{date}", app.GetByDate)
	mux.Post("/exercises", app.AddNewExercise)
	mux.Patch("/exercises/{date}/{name}/{count}", app.UpdateExercise)
	mux.Delete("/exercises/{date}/{name}", app.DeleteExercise)

	return mux
}
