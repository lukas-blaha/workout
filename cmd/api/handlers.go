package main

import (
	"net/http"
	"time"
)

var Exercises = []jsonResponse{
	{Name: "pushups", Count: 0, Date: getDate(time.Now())},
	{Name: "pullups", Count: 0, Date: getDate(time.Now())},
	{Name: "dips", Count: 0, Date: getDate(time.Now())},
	{Name: "squats", Count: 0, Date: getDate(time.Now())},
}

func (app *Config) GetExercises(w http.ResponseWriter, r *http.Request) {
	for _, e := range Exercises {
		app.writeJSON(w, http.StatusOK, e)
	}
}
