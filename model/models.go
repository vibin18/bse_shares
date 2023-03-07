package model

import "time"

type ShareReport struct {
	Name  string  `json:"name"`
	Count int     `json:"count"`
	Total float32 `json:"total"`
}

type Share struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SellBuyShare struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Count     int       `json:"count"`
	Price     float32   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Type      string    `json:"type"`
}

type TotalShare struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	PCount int    `json:"pc"`
	SCount int    `json:"sc"`
	TCount int    `json:"tc"`
}
type SellShare struct {
	Name      string  `json:"name"`
	Count     int     `json:"count"`
	Price     float32 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

type ShareReports []ShareReport
