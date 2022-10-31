package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (app *Config) GetByNameDate(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	date := chi.URLParam(r, "date")

	exercise, err := app.Models.Exercise.GetByNameDate(name, date)
	if err != nil {
		log.Panic("Can't return exercises!")
	}

	app.writeJSON(w, http.StatusOK, exercise)
}

func (app *Config) GetByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	exercise, err := app.Models.Exercise.GetByName(name)
	if err != nil {
		log.Panic("Can't return exercises!")
	}

	app.writeJSON(w, http.StatusOK, exercise)
}

func (app *Config) GetByDate(w http.ResponseWriter, r *http.Request) {
	date := chi.URLParam(r, "date")

	exercise, err := app.Models.Exercise.GetByDate(date)
	if err != nil {
		log.Panic("Can't return exercises!")
	}

	app.writeJSON(w, http.StatusOK, exercise)
}

func (app *Config) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	date := chi.URLParam(r, "date")
	c := chi.URLParam(r, "count")

	count, err := strconv.Atoi(c)
	if err != nil {
		app.errorJSON(w, errors.New("Can't convert count to integer."))
		return
	}

	err = app.Models.Exercise.Update(name, date, count)
}

func (app *Config) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	date := chi.URLParam(r, "date")

	app.Models.Exercise.Delete(name, date)
}
