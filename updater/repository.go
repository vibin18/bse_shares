package updater

import "github.com/vibin18/bse_shares/model"

type ShareListUpdaterRepository interface {
	GetAllPurchaseReport() (model.ShareReports, error)
}

type ShareCacheUpdaterRepository interface {
}
