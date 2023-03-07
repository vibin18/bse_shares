package config

import (
	"github.com/vibin18/bse_shares/updater"
)

type AppConfig struct {
	Data      []*updater.Stock
	ShareList []*string
}

type ShareReport struct {
	Name  string  `json:"name"`
	Count int     `json:"count"`
	Total float32 `json:"total"`
}
