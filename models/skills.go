package models

type UserSkills struct {
	ID         string
	Name       string
	Created_at string
	Updated_at string
}

// UpdateSkills => data format to update skills of a user
type UpdateSkills struct {
	Skills string
}
