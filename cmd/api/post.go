package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Nostafa/golang-jwt/internal/store"
	"github.com/go-chi/chi/v5"
)

type createPostRequest struct {
	Content string   `json:"content" validate:"required,max=500"`
	Title   string   `json:"title" validate:"required,max=50"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input createPostRequest
	if err := readJson(w, r, &input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	userId := 1
	ctx := r.Context()
	post := store.Post{
		Content: input.Content,
		Title:   input.Title,
		Tags:    input.Tags,
		UserId:  int64(userId),
	}

	err := app.store.Post.Create(ctx, &post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = writeJson(w, http.StatusCreated, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
func (app *application) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	// userId := 1

	ctx := r.Context()
	postId := chi.URLParam(r, "postId")
	postIdInt, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post, err := app.store.Post.GetById(ctx, postIdInt)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = writeJson(w, http.StatusOK, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
