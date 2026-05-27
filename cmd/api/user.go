package main

import (
	"net/http"

	"github.com/Nostafa/golang-jwt/internal/store"
)

type createUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input createUserRequest
	if err := readJson(w, r, &input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Email:     input.Email,
		Password:  []byte(input.Password),
	}
	ctx := r.Context()

	err := app.store.User.Create(ctx, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = writeJson(w, http.StatusCreated, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
