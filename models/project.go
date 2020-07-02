package models

import (
	"database/sql"
)

type OwnerInfo struct {
	ID               int    `json:"id"`
	Email            string `json:"email"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Location         string `json:"location"`
	Member           string `json:"member"`
	ProjectCompleted int    `json:"projectCompleted"`
	PhoneCode        int    `json:"phoneCode"`
	PhoneNumber      int    `json:"phoneNumber"`
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
	ID           int    `json:"id"`
	Project_link string `json:"projectLink"`
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

type UserReviewInfo struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Picture   string `json:"picture"`
}

type GetUserProjectResponse struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	IsOwner     bool           `json:"isOwner"`
	TrelloURL   string         `json:"trelloUrl"`
	IsComment   bool           `json:"isComment"`
	Reviewee    UserReviewInfo `json:"reviewee"`
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
	Status            string
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
	Status            string                          `json:"status"`
}

type ProjectDetailInterestedMember struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Picture   string `json:"picture"`
}

type ProjectAcceptMemberParameter struct {
	FreelancerID int    `json:"freelancerID"`
	TrelloKey    string `json:"trelloKey"`
}
