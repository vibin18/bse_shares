package handler

import (
	"github.com/vibin18/bse_shares/model"
)

type HandlerRepo interface {
	GetAllPurchaseReport() (model.ShareReports, error)
	GetAllSalesReport() ([]model.ShareReport, error)
	GetAllShares() ([]model.Share, error)
	InsertNewShare(res model.Share) error
	GetShareByID(id int) (model.Share, error)
	BuyShare(res model.SellBuyShare) error
	SellShare(res model.SellBuyShare) error
	GetAllSharesWithData() ([]model.TotalShare, error)
	GetAllPurchases() ([]model.SellBuyShare, error)
	GetAllSales() ([]model.SellBuyShare, error)
}
