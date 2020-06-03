package models

type ExperienceParameters struct {
	Place     string
	Position  string
	StartYear int
	EndYear   int
}

type ExperienceReturnValue struct {
	ID int
	ExperienceParameters
}

type ExperienceTableResponse struct {
	ExperienceReturnValue
	UserID int
}

type AddExperienceParameter struct {
	Experience []ExperienceParameters
}
