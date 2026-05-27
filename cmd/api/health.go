package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":      "ok",
		"environment": app.config.env,
		"version":     app.config.version,
	}

	err := writeJson(w, http.StatusOK, data)
	if err != nil {
		errorJson(w, http.StatusInternalServerError, err.Error())
	}
}
