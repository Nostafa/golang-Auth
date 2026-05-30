package main

import (
	"net/http"
	"strconv"

	"github.com/Nostafa/golang-jwt/internal/store"
	"github.com/go-chi/chi/v5"
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
		Password:  input.Password,
	}
	ctx := r.Context()

	_, err := app.store.User.Create(ctx, user)
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

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := chi.URLParam(r, "userId")
	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user, err := app.store.User.GetById(ctx, userIdInt)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = writeJson(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
