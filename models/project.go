package models

type SearchProject struct {
}

type ProjectResult struct {
	ID         string
	Title      string
	Skills     string
	Price      float64
	Attachment string
}

type ProjectSearchResponse struct {
	ID         string
	Title      string
	Skills     []string
	Price      float64
	Attachment []string
}

type ProjectLinksResponse struct {
	ID           string
	Project_link string
}
