package helpers

import (
	"PMSFreelancer/config"
	"errors"
)

func SkillList(s string) error {
	splittedData := SplitComma(s)

	var allSkills []interface{}

	rows, err := config.DB.Query("SELECT id from skills")

	if err != nil {
		return err
	}

	var skillid string
	for rows.Next() {
		err := rows.Scan(&skillid)

		if err != nil {
			return err
		}

		allSkills = append(allSkills, skillid)
	}

	for i := 0; i < len(splittedData); i++ {
		exist := Contains(allSkills, splittedData[i])

		if !exist {
			return errors.New("not exist")
		}
	}

	return nil
}
