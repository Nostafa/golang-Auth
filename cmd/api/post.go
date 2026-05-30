package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Nostafa/golang-jwt/internal/store"
	"github.com/go-chi/chi/v5"
)

type postCtxKey string

const postKey postCtxKey = "post"

type createPostRequest struct {
	Content string   `json:"content" validate:"required,max=500"`
	Title   string   `json:"title" validate:"required,max=50"`
	Tags    []string `json:"tags"`
}

type updatePostRequest struct {
	Content *string   `json:"content" validate:"omitempty,max=500"`
	Title   *string   `json:"title" validate:"omitempty,max=50"`
	Tags    *[]string `json:"tags" validate:"omitempty"`
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

	_, err := app.store.Post.Create(ctx, &post)
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

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post, err := getPostFromContext(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	var input updatePostRequest
	if err := readJson(w, r, &input); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Tags != nil {
		post.Tags = *input.Tags
	}

	updatedPost, err := app.store.Post.Update(r.Context(), post.Id, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	err = writeJson(w, http.StatusOK, updatedPost)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId := chi.URLParam(r, "postId")
		postIdInt, err := strconv.ParseInt(postId, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()
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
		ctx = context.WithValue(ctx, postKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(r *http.Request) (*store.Post, error) {
	post, ok := r.Context().Value(postKey).(*store.Post)
	if !ok {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postId := chi.URLParam(r, "postId")
	postIdInt, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err = app.store.Post.Delete(ctx, postIdInt)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = writeJson(w, http.StatusOK, nil)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
