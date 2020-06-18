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
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type SearchProjectResponse struct {
	Project     []SearchProjectQuery
	TotalSearch int `json:"totalSearch"`
}

type ProjectLinksResponse struct {
	ID           int
	Project_link string
}

type CreateProject struct {
	Title       string
	Description string
	Skills      []string
	Price       float64
	Attachment  []string
}

type GetUserProjectResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type ProjectDetailRequest struct {
	ID                int
	Title             string
	Skills            string
	Price             float64
	OwnerID           int
	InterestedMembers sql.NullString
}

type ProjectDetailResponse struct {
	ID                int                             `json:"id"`
	Title             string                          `json:"title"`
	Skills            []string                        `json:"skills"`
	Attachment        []string                        `json:"attachment"`
	Price             float64                         `json:"price"`
	Owner             OwnerInfo                       `json:"owner"`
	InterestedMembers []ProjectDetailInterestedMember `json:"interestedMembers"`
}

type ProjectDetailInterestedMember struct {
	ID        int
	FirstName string
	LastName  string
}

type ProjectAcceptMemberParameter struct {
	FreelancerID int    `json:"freelancerID"`
	TrelloKey    string `json:"trelloKey"`
}
