package updater

import "github.com/vibin18/bse_shares/utils"

type ListUpdaterService interface {
	ListUpdate() []*string
}

type CacheUpdaterService interface {
	Update([]*string) []*utils.Stock
}
