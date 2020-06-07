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

type GoogleResponse struct {
	ID             string
	Email          string
	Verified_email bool
	Picture        string
}

type UpdateResetPassword struct {
	Password string
	Token    string
}

type DatabaseResetPassword struct {
	Email  string
	Token  string
	Expire string
}

type UserProfile struct {
	ID          string
	Fullname    string
	Email       string
	Description string
	Education   []EducationReturnValue
	Skill       []UserSkills
	Experience  []ExperienceReturnValue
	Picture     string
	Username    string
	Location    string
	Member      string
}

type QueryUserProfile struct {
	ID          string
	Firstname   string
	LastName    string
	Email       string
	Description string
	Picture     string
	CreatedAt   string
	Username    string
	Location    string
	Skills      string
}
