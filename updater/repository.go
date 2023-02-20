package updater

type ShareListUpdaterRepository interface {
	GetAllPurchaseReport() (ShareReports, error)
}

type ShareCacheUpdaterRepository interface {
}
