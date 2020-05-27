package models

// RegistrationUserUsingPassword => User Model for registration via email and password
type RegistrationUserUsingPassword struct {
	Email       string
	Password    string
	Fullname    string
	Location    int
	Description string
	Skills      string
	Username    string
}

type LoginUserPassword struct {
	Email    string
	Password string
}
