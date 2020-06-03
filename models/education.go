package models

type EducationParameters struct {
	Name      string
	StartYear int
	EndYear   int
}

type EducationReturnValue struct {
	ID int
	EducationParameters
}

type EducationTableResponse struct {
	EducationReturnValue
	UserID int
}

type AddEducationParameter struct {
	Education []EducationParameters
}
