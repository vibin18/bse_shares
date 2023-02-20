package updater

import (
	"log"
)

type listUpdaterService struct {
	listUpdaterRepo ShareListUpdaterRepository
}

func NewListUpdaterService(listUpdaterRepo ShareListUpdaterRepository) ListUpdaterService {
	return &listUpdaterService{
		listUpdaterRepo: listUpdaterRepo,
	}
}

func (s *listUpdaterService) ListUpdate() []string {
	log.Println("Updating share list")
	var ul []string

	log.Println("Fetching db for new list")
	dl, err := s.listUpdaterRepo.GetAllPurchaseReport()

	if err != nil {
		log.Println("Error getting share list report from postgres, error : ", err)
	}
	for _, share := range dl {
		ul = append(ul, share.Name)
	}
	return ul
}
