package models

import "database/sql"

type OwnerInfoQuery struct {
	ID        string
	FirstName string
	LastName  string
	Location  string
	CreatedAt string
}

type OwnerInfo struct {
	ID               string
	FullName         string
	Location         string
	Member           string
	ProjectCompleted int
}

type SearchProjectQuery struct {
	ID          string
	Title       string
	Description string
	Price       float64
}

type SearchProjectResponse struct {
	Project     []SearchProjectQuery
	TotalSearch int
}

type ProjectLinksResponse struct {
	ID           string
	Project_link string
}

type CreateProject struct {
	Title       string
	Description string
	Skills      string
	Price       float64
	Attachment  string
}

type GetUserProjectResponse struct {
	ID          string
	Title       string
	Description string
	Status      string
}

type ProjectDetailRequest struct {
	ID                string
	Title             string
	Skills            string
	Price             float64
	Attachment        sql.NullString
	OwnerID           int
	InterestedMembers sql.NullString
}

type ProjectDetailResponse struct {
	ID         string
	Title      string
	Skills     []string
	Attachment []string
	Price      float64
	Owner      OwnerInfo
}
