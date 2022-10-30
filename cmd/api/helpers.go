package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Exercise struct {
	ID      int    `json:"id,omitepmpty"`
	Name    string `json:"name,omitempty"`
	Count   int    `json:"count,omitempty"`
	Date    string `json:"date,omitempty"`
	Error   bool   `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

var Exercises = []Exercise{
	{Name: "pushups", Count: 0, Date: "30-10-2022"},
	{Name: "pullups", Count: 0, Date: "30-10-2022"},
	{Name: "dips", Count: 0, Date: "30-10-2022"},
	{Name: "squats", Count: 0, Date: "30-10-2022"},
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if status != http.StatusOK && status != http.StatusAccepted {
		w.WriteHeader(status)
	}

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload Exercise
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

func getDate(t time.Time) string {
	return fmt.Sprintf("%02d-%02d-%d", t.Day(), t.Month(), t.Year())
}
