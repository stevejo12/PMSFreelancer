package models

type ResponseWithNoBody struct {
	code    string
	message string
}

type ResponseWithStringData struct {
	ResponseWithNoBody
	data string
}

type ResponseOKLogin struct {
	ResponseWithNoBody
	data TokenResponse
}

type ResponseOKGetAllCountries struct {
	ResponseWithNoBody
	data []CountryData
}

type ResponseOKGetUserProfile struct {
	ResponseWithNoBody
	data UserProfile
}

type ResponseOKGetUserEducation struct {
	ResponseWithNoBody
	data []EducationReturnValue
}

type ResponseOKGetUserExperience struct {
	ResponseWithNoBody
	data []ExperienceReturnValue
}

type ResponseOKGetUserProject struct {
	ResponseWithNoBody
	data []GetUserProjectResponse
}

type ResponseOKProjectDetail struct {
	ResponseWithNoBody
	data ProjectDetailResponse
}

type ParamSearchProject struct {
	page int
	size int
}

type ResponseOKUserReviews struct {
	ResponseWithNoBody
	data []ReviewInfo
}

type ParamFilterProject struct {
	page int
	size int
}

type ResponseOKSearchProject struct {
	ResponseWithNoBody
	data []SearchProjectResponse
}

type ResponseOKProjectCategory struct {
	ResponseWithNoBody
	data []ProjectCategoryData
}
