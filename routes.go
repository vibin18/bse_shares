package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vibin18/bse_shares/handler"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.DefaultLogger)

	mux.Get("/", handler.IndexHandler)
	mux.Post("/create", handler.PostUpdateHandler)

	return mux
}
