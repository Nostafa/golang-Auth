package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %v \n path: %s \n error: %s ", err, r.URL.Path, err.Error())
	errorJson(w, http.StatusInternalServerError, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %v \n path: %s \n error: %s ", err, r.URL.Path, err.Error())
	errorJson(w, http.StatusNotFound, err.Error())
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %v \n path: %s \n error: %s ", err, r.URL.Path, err.Error())
	errorJson(w, http.StatusBadRequest, err.Error())
}
