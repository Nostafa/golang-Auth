package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Nostafa/golang-jwt/internal/store"
	"github.com/go-chi/chi/v5"
)

type createCommentRequest struct {
	Content string `json:"content" validate:"required,max=50"`
	PostId  int64  `json:"postId" validate:"required,min=1"`
	UserId  int64  `json:"userId" validate:"required,min=1"`
}

type createCommentResponse struct {
	Id        int64     `json:"id"`
	Content   string    `json:"content"`
	PostId    int64     `json:"postId"`
	UserId    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input createCommentRequest
	if err := readJson(w, r, &input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comment := store.Comment{
		Content: input.Content,
		PostId:  input.PostId,
		UserId:  input.UserId,
	}
	ctx := r.Context()
	err := app.store.Comment.Create(ctx, &comment)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	response := createCommentResponse{
		Id:        comment.Id,
		Content:   input.Content,
		PostId:    input.PostId,
		UserId:    input.UserId,
		CreatedAt: comment.CreatedAt,
	}
	err = writeJson(w, http.StatusCreated, response)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getCommentsByPostIdHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	postId := chi.URLParam(r, "postId")
	postIdInt, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comments, err := app.store.Comment.GetByPostId(ctx, postIdInt)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = writeJson(w, http.StatusOK, comments)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
