package models

type PortfolioDatabase struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	OwnerID     string `json:"ownerId"`
}

type PortfolioEditParameter struct {
	Description string
}
