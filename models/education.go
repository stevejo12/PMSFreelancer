package models

type EducationParameters struct {
	Name      string `json:"name"`
	StartYear int    `json:"startYear"`
	EndYear   int    `json:"endYear"`
}

type EducationReturnValue struct {
	ID int `json:"id"`
	EducationParameters
}

type EducationTableResponse struct {
	EducationReturnValue
	UserID int
}
