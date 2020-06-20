package models

type ReviewDatabase struct {
	ID int
	AddReview
	MemberID  int
	CreatedAt string
}

type AddReview struct {
	Message    string
	StarRating int
	ProjectID  int
	UserID     int
}

type ProjectInformationForReview struct {
	ID    int
	Title string
}

type UserInformationForReview struct {
	ID        int
	FirstName string
	LastName  string
}

type ReviewInfo struct {
	ID         int
	Message    string
	StarRating int
	Project    ProjectInformationForReview
	Reviewer   UserInformationForReview
	CreatedAt  string
}
