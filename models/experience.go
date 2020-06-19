package models

type ExperienceParameters struct {
	Place       string `json:"place"`
	Position    string `json:"position"`
	StartYear   int    `json:"startYear"`
	EndYear     int    `json:"endYear"`
	Description string `json:"description"`
}

type ExperienceReturnValue struct {
	ID int `json:"id"`
	ExperienceParameters
}

type ExperienceTableResponse struct {
	ExperienceReturnValue
	UserID int
}

type AddExperienceParameter struct {
	Experience []ExperienceParameters
}
