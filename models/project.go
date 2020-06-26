package models

import (
	"database/sql"
)

type OwnerInfoQuery struct {
	ID        string
	FirstName string
	LastName  string
	Location  string
	CreatedAt string
}

type OwnerInfo struct {
	ID               string `json:"id"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Location         string `json:"location"`
	Member           string `json:"member"`
	ProjectCompleted int    `json:"projectCompleted"`
}

type SearchProjectQuery struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Skill       []string `json:"skill"`
	CreatedAt   string   `json:"createdAt"`
	Category    string   `json:"category"`
}

type SearchProjectResponse struct {
	Project  []SearchProjectQuery `json:"project"`
	PageMeta PageInfoSchema       `json:"PageMeta"`
}

type PageInfoSchema struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type ProjectLinksResponse struct {
	ID           int
	Project_link string
}

type CreateProject struct {
	Title       string
	Description string
	Skills      []int
	Price       float64
	Attachment  []string
	Category    int
}

type EditProject struct {
	Title       string
	Description string
	Skills      []int
	Price       float64
	Attachment  []string
	Category    int
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
	Description       string
	Category          int
}

type ProjectDetailResponse struct {
	ID                int                             `json:"id"`
	Title             string                          `json:"title"`
	Skills            []string                        `json:"skills"`
	Attachment        []ProjectLinksResponse          `json:"attachment"`
	Price             float64                         `json:"price"`
	Owner             OwnerInfo                       `json:"owner"`
	InterestedMembers []ProjectDetailInterestedMember `json:"interestedMembers"`
	Description       string                          `json:"description"`
	Category          string                          `json:"category"`
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
