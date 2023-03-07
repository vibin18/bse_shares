package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/vibin18/bse_shares/config"
	"github.com/vibin18/bse_shares/handler"
	"github.com/vibin18/bse_shares/opts"
	"github.com/vibin18/bse_shares/repository/postgres"
	"github.com/vibin18/bse_shares/updater"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	argparser *flags.Parser
	arg       opts.Params
)

func initArgparser() {
	argparser = flags.NewParser(&arg, flags.Default)
	_, err := argparser.Parse()

	// check if there is a parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

func main() {
	initArgparser()
	var app config.AppConfig
	log.Println("Connecting to DB...")
	dsn := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v", arg.DbServer, arg.DbPort, arg.DbName, arg.DbUser, arg.DbPass)
	db, err := postgres.NewPgSQL(dsn)
	if err != nil {
		log.Fatal(err)
	}

	myShares := updater.NewListUpdaterService(db)
	app.ShareList = myShares.ListUpdate()
	handler.CreateNewHandlerConfig(&app)

	go func() {
		for range time.Tick(time.Second * 10) {
			myShareCache := updater.NewCacheUpdaterService(db)
			myStockList := myShareCache.Update(app.ShareList)
			mutex := sync.Mutex{}
			mutex.Lock()
			defer mutex.Unlock()
			app.Data = myStockList
		}
	}()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(db),
	}
	log.Printf("Starting HTTP server")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
