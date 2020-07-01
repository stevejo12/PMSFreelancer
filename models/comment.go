package models

type AddReview struct {
	Message    string
	StarRating int
	ProjectID  int
	IsOwner    bool
}

type ProjectInformationForReview struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type UserInformationForReview struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsOwner   bool   `json:"isOwner"`
}

type ReviewInfo struct {
	ID         int                         `json:"id"`
	Message    string                      `json:"message"`
	StarRating int                         `json:"starRating"`
	Project    ProjectInformationForReview `json:"project"`
	Reviewer   UserInformationForReview    `json:"reviewer"`
	CreatedAt  string                      `json:"createdAt"`
}
