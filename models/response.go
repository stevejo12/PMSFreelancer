package models

type ResponseWithNoBody struct {
	code    string
	message string
}

type ParamSearchProject struct {
	page int
	size int
}

type ParamFilterProject struct {
	page int
	size int
}
