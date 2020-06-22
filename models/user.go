package models

// RegistrationUserUsingPassword => User Model for registration via email and password
type RegistrationUserUsingPassword struct {
	Email       string
	Password    string
	FirstName   string
	LastName    string
	Location    int
	Description string
	Skills      []int
	Username    string
}

type RegistrationUserUsingGoogle struct {
	Email       string
	GoogleID    string
	FirstName   string
	LastName    string
	Location    int
	Description string
	Skills      []int
	Username    string
}

// LoginUserPassword => Parameter login
type LoginUserPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GoogleResponse struct {
	ID             string
	Email          string
	Verified_email bool
	Picture        string
}

type RegistrationInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
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
	ID               int                        `json:"id"`
	FirstName        string                     `json:"firstName"`
	LastName         string                     `json:"lastName"`
	Email            string                     `json:"email"`
	Description      string                     `json:"description"`
	Education        []EducationReturnValue     `json:"education"`
	Skill            []UserSkills               `json:"skill"`
	Experience       []ExperienceReturnValue    `json:"experience"`
	Picture          string                     `json:"picture"`
	Username         string                     `json:"username"`
	Location         CountryDataProfile         `json:"location"`
	Member           string                     `json:"member"`
	Portfolio        []PortfolioReturnParameter `json:"portfolio"`
	Balance          float64                    `json:"balance"`
	ProjectCompleted int                        `json:"projectCompleted"`
}

type QueryUserProfile struct {
	ID          int
	Firstname   string
	LastName    string
	Email       string
	Description string
	Picture     string
	CreatedAt   string
	Username    string
	Location    int
	Skills      string
	Balance     float64
}

type UpdateProfile struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Username    string `json:"username"`
	Location    int    `json:"location"`
	Description string `json:"description"`
}
