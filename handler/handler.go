package handler

import (
	"encoding/json"
	"github.com/vibin18/bse_shares/config"
	"net/http"
)

var app *config.AppConfig

func CreateNewHandlerConfig(a *config.AppConfig) {
	app = a
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.Data)
}

//func PostUpdateHandler(w http.ResponseWriter, r *http.Request) {
//	var stock *string
//	s := r.FormValue("stock")
//	stock = &s
//	app.ShareList = append(app.ShareList, stock)
//}
