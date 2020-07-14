package helpers

import (
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"

	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
	"errors"
)

// GetSkillRaw => get id and name from skill id
func GetSkillRaw(param string) ([]models.SkillRaw, error) {
	var skillRaw []models.SkillRaw
	initialQuery, err := SettingInQueryWithID("skills", param, "id, name")

	if err != nil {
		return skillRaw, err
	}

	data, err := config.DB.Query(initialQuery)
	defer data.Close()

	if err != nil {
		return skillRaw, errors.New("Server is unable to execute query to the database")
	}

	for data.Next() {
		var skills models.SkillRaw
		if err := data.Scan(&skills.ID, &skills.Name); err != nil {
			return []models.SkillRaw{}, errors.New("Something is wrong with the database data")
		}

		skillRaw = append(skillRaw, skills)
	}

	return skillRaw, nil
}
