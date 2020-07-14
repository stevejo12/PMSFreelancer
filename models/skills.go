package models

type UserSkills struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Created_at string `json:"createdAt"`
	Updated_at string `json:"updatedAt"`
}

// UpdateSkills => data format to update skills of a user
type UpdateSkills struct {
	Skills []int
}

type SkillRaw struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
