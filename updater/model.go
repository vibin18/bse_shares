package updater

type ShareReport struct {
	Name  string  `json:"name"`
	Count int     `json:"count"`
	Total float32 `json:"total"`
}

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

type ShareReports []ShareReport
