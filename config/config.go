package config

import (
	"github.com/vibin18/bse_shares/updater"
)

type AppConfig struct {
	Data      []*updater.Stock
	ShareList []*string
}
