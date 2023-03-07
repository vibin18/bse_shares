package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vibin18/bse_shares/handler"
	"net/http"
)

func routes(db handler.HandlerRepo) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.DefaultLogger)

	handlers := handler.NewHandlerService(db)

	mux.Get("/", handler.IndexHandler)
	mux.Post("/add-share-post", handlers.AddSharesPostHandler)
	mux.Get("/update-share", handlers.UpdateShareHandler)
	mux.Post("/update-share-post", handlers.UpdateSharePostHandler)
	mux.Get("/error", handler.ErrorHandler)
	mux.Get("/add-share", handler.AddSharesHandler)
	mux.Get("/get-sales", handlers.ListAllSalesHandler)
	mux.Get("/get-purchases", handlers.ListAllPurchaseHandler)
	mux.Get("/get-purchase-report", handlers.ReportAllPurchaseHandler)
	mux.Get("/get-sales-report", handlers.ReportAllSalesHandler)
	mux.Get("/list-share", handlers.ListTotalSharesHandler)
	mux.Get("/stock", handlers.StockHandler)

	return mux
}
