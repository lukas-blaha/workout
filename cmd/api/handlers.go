package main

import (
	"log"
	"net/http"

	"github.com/lukas-blaha/workout/data"
)

func (app *Config) GetExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := app.Models.Exercise.GetAll()
	if err != nil {
		log.Panic("Can't return exercises!")
	}

	for _, e := range exercises {
		app.writeJSON(w, http.StatusOK, e)
	}
}

func (app *Config) AddNewExercise(w http.ResponseWriter, r *http.Request) {
	var exercise data.Exercise
	err := app.readJSON(w, r, &exercise)
	if err != nil {
		log.Fatal(err)
	}

	app.Models.Exercise = exercise
	app.Models.Exercise.Insert()

	err = app.writeJSON(w, http.StatusOK, exercise)
	if err != nil {
		log.Fatal(err)
	}
}
