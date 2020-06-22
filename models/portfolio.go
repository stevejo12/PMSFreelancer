package models

type PortfolioRequestParameter struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	StartYear   int    `json:"startYear"`
	EndYear     int    `json:"endYear"`
}

type PortfolioReturnParameter struct {
	ID int `json:"id"`
	PortfolioRequestParameter
}

type PortfolioEditParameter struct {
	Description string
}
