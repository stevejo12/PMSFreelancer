package models

import "mime/multipart"

type PortfolioRequestParameter struct {
	File multipart.File
}

type PortfolioDatabase struct {
	ID          string
	Description string
	Link        string
	OwnerID     string
}
