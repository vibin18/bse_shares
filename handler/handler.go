package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vibin18/bse_shares/config"
	"github.com/vibin18/bse_shares/model"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var app *config.AppConfig
var tmpl *template.Template

func CreateNewHandlerConfig(a *config.AppConfig) {
	app = a
}

//func PostUpdateHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	if r.Method != "POST" {
//		w.Write(bytes.NewBufferString("Invalid request").Bytes())
//	}
//	if r.FormValue("stock") == "" {
//		w.Write(bytes.NewBufferString("Invalid request, Wrong data entered!").Bytes())
//		return
//	}
//	var stock *string
//
//	s := r.FormValue("stock")
//
//	stock = &s
//	app.ShareList = append(app.ShareList, stock)
//	log.Println("Updating new list")
//	w.Write(bytes.NewBufferString("Stock added").Bytes())
//}

type handlerService struct {
	HandlerRepo HandlerRepo
}

func NewHandlerService(HandlerRepo HandlerRepo) handlerService {
	return handlerService{
		HandlerRepo: HandlerRepo,
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "index.html", app.Data)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	buff.WriteTo(w)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	tmpl.ExecuteTemplate(w, "error_main.html", nil)
}

func AddSharesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	tmpl.ExecuteTemplate(w, "add_shares_main.html", nil)
}

func (d *handlerService) ListSharesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []model.Share{}

	shares, err = d.HandlerRepo.GetAllShares()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	tmpl.ExecuteTemplate(w, "list_shares_main.html", shares)
}

func random(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (d *handlerService) AddSharesPostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to parse template" + err.Error())
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	name := r.Form.Get("company")

	id := random(10, 9999999)

	// log.Printf("Adding ID %v", id)

	myShare := model.Share{
		Name: name,
		Id:   id,
	}
	err = d.HandlerRepo.InsertNewShare(myShare)
	if err != nil {
		log.Println("DB execution failed to add shares " + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	buff := bytes.Buffer{}

	log.Println(myShare.Name + " added")

	tmpl.ExecuteTemplate(&buff, "add_shares_success_main.html", myShare)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	buff.WriteTo(w)
}

func (d *handlerService) UpdateShareHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to parse template" + err.Error())
	}

	shares := []model.Share{}
	shares, err = d.HandlerRepo.GetAllShares()

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "update_shares_main.html", shares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	buff.WriteTo(w)
}

func (d *handlerService) UpdateSharePostHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	name := r.Form.Get("name")
	cdate := r.Form.Get("date")
	datetime, err := time.Parse("2006-01-02", cdate)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	c := r.Form.Get("count")
	count, err := strconv.Atoi(c)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	log.Printf("Date entered : %v", cdate)
	shareIdName := strings.Split(name, "---")

	idn, err := strconv.Atoi(shareIdName[1])
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	p := r.Form.Get("price")
	price, err := strconv.ParseFloat(p, 32)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	pf := float32(price)
	um := r.Form.Get("update_type")
	var umType string

	buySellShare := model.SellBuyShare{
		Id:        idn,
		Name:      shareIdName[0],
		Count:     count,
		Price:     pf,
		CreatedAt: datetime,
		UpdatedAt: datetime,
		Type:      umType,
	}

	if um == "Buy" {
		umType = "bought"
		buySellShare.Type = umType
		log.Printf("%v %v shares of %v", um, count, shareIdName[0])
		err = d.HandlerRepo.BuyShare(buySellShare)
		if err != nil {
			log.Println("DB execution failed to add shares " + err.Error())
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
	} else {
		umType = "sold"
		buySellShare.Type = umType
		log.Printf("%v %v shares of %v", um, count, shareIdName[0])
		err = d.HandlerRepo.SellShare(buySellShare)
		if err != nil {
			log.Println("DB execution failed to add shares " + err.Error())
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
	}

	log.Printf("%v %v shares of %v", um, count, shareIdName[0])

	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to parse template" + err.Error())
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "update_shares_success_main.html", buySellShare)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	buff.WriteTo(w)
}

func (d *handlerService) ListTotalSharesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []model.TotalShare{}
	ushares := []model.TotalShare{}

	shares, err = d.HandlerRepo.GetAllSharesWithData()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	for _, share := range shares {
		share.TCount = share.PCount - share.SCount
		ushares = append(ushares, share)
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "list_shares_main.html", ushares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	buff.WriteTo(w)
}

func (d *handlerService) ListAllPurchaseHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []model.SellBuyShare{}
	pshares := []model.SellShare{}

	shares, err = d.HandlerRepo.GetAllPurchases()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	for _, share := range shares {
		// share.CreatedAt.Format("02-Jan-2006")

		p := model.SellShare{
			CreatedAt: share.CreatedAt.Format("2006-Jan-02"),
			Name:      share.Name,
			Count:     share.Count,
			Price:     share.Price,
		}
		pshares = append(pshares, p)
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "list_purchases_main.html", pshares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	buff.WriteTo(w)
}

func (d *handlerService) ListAllSalesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []model.SellBuyShare{}
	pshares := []model.SellShare{}

	shares, err = d.HandlerRepo.GetAllSales()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	for _, share := range shares {
		share.CreatedAt.Format("02-Jan-2006")

		p := model.SellShare{
			CreatedAt: share.CreatedAt.Format("02-Jan-2006"),
			Name:      share.Name,
			Count:     share.Count,
			Price:     share.Price,
		}
		pshares = append(pshares, p)
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "list_sales_main.html", pshares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	buff.WriteTo(w)
}

func (d *handlerService) ReportAllPurchaseHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	var shares []model.ShareReport

	shares, err = d.HandlerRepo.GetAllPurchaseReport()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "list_purchase_report_main.html", shares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	buff.WriteTo(w)
}

func (d *handlerService) ReportAllSalesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Println("Failed to pasre template" + err.Error())
	}
	shares := []model.ShareReport{}

	shares, err = d.HandlerRepo.GetAllSalesReport()
	if err != nil {
		log.Println("Failed list shares from DB" + err.Error())
	}

	buff := bytes.Buffer{}

	tmpl.ExecuteTemplate(&buff, "list_sales_report_main.html", shares)
	if err != nil {
		log.Println("Failed to execute template" + err.Error())
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	buff.WriteTo(w)
}

func (d *handlerService) StockHandler(w http.ResponseWriter, r *http.Request) {
	data := app.Data
	json, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed marshal json" + err.Error())
	}
	io.Writer(w).Write(json)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Share added.")
}
