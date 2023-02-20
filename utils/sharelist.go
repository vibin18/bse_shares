package utils

type ShareCommon struct {
	ShareDataList *[]Stock
}

//func NewShareCommon(shareDataList []*Stock) *ShareCommon {
//	return &ShareCommon{ShareDataList: shareDataList}
//}
//
//func (sc *ShareCommon) SetShareCommon(s []*Stock) {
//	sc.ShareDataList = s
//}
//
//func (sc *ShareCommon) GetShareCommon() []*Stock {
//	return sc.ShareDataList
//}

type Stock struct {
	CompanyName   string
	ShortName     string
	CurrentValue  string
	PreviousClose string
	PreviousOpen  string
	DayHigh       string
	DayLow        string
	WeekAverage   string
	UpdateTime    string
}
