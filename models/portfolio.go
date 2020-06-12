package models

import "mime/multipart"

type PortfolioRequestParameter struct {
	File multipart.File
}

type PortfolioDatabase struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Link        string `json:"link"`
	OwnerID     string `json:"ownerId"`
}
