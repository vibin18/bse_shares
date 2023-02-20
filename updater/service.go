package updater

type ListUpdaterService interface {
	ListUpdate() []*string
}

type CacheUpdaterService interface {
	Update([]*string) []*Stock
}
