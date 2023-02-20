package config

import "github.com/vibin18/bse_shares/utils"

type AppConfig struct {
	Data      []*utils.Stock
	ShareList []*string
}
