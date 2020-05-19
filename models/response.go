package models

type ResponseWithNoBody struct {
	code    string
	message string
}

type ResponseLoginWithToken struct {
	ResponseWithNoBody
	token string
}
