package handler

import "net/http"

type HandlerService interface {
	ListSharesHandler(w http.ResponseWriter, r *http.Request)
	AddSharesPostHandler(w http.ResponseWriter, r *http.Request)
	UpdateShareHandler(w http.ResponseWriter, r *http.Request)
	UpdateSharePostHandler(w http.ResponseWriter, r *http.Request)
	ListTotalSharesHandler(w http.ResponseWriter, r *http.Request)
	ListAllPurchaseHandler(w http.ResponseWriter, r *http.Request)
	ListAllSalesHandler(w http.ResponseWriter, r *http.Request)
	ReportAllPurchaseHandler(w http.ResponseWriter, r *http.Request)
	ReportAllSalesHandler(w http.ResponseWriter, r *http.Request)
	StockHandler(w http.ResponseWriter, r *http.Request)
}
