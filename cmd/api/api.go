package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Nostafa/golang-jwt/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  *store.Storage
}

type config struct {
	port    string
	db      dbConfig
	env     string
	version string
}

type dbConfig struct {
	url          string
	maxOpenConns int
	maxIdleConns int
	maxLifetime  string
	maxIdleTime  string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr) // pick one ClientIPFrom* based on your infra, see below
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", app.healthCheckHandler)
	r.Route("/posts", func(r chi.Router) {
		r.Post("/", app.createPostHandler)

		r.Route("/{postId}", func(r chi.Router) {
			r.Get("/", app.getPostByIdHandler)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", app.createUserHandler)

	})

	return r
}

func (app *application) server(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.port,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("🚀 Starting server on port %s", app.config.port)

	return srv.ListenAndServe()
}
