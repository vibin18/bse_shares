package handler

import (
	"bytes"
	"encoding/json"
	"github.com/vibin18/bse_shares/config"
	"log"
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

func PostUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.Write(bytes.NewBufferString("Invalid request").Bytes())
	}
	if r.FormValue("stock") == "" {
		w.Write(bytes.NewBufferString("Invalid request, Wrong data entered!").Bytes())
		return
	}
	var stock *string

	s := r.FormValue("stock")

	stock = &s
	app.ShareList = append(app.ShareList, stock)
	log.Println("Updating new list")
	w.Write(bytes.NewBufferString("Stock added").Bytes())
}
