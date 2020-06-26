package helpers

import (
	"github.com/stevejo12/PMSFreelancer/config"
	// "PMSFreelancer/config"
	"errors"
)

func GetCategoryNameByID(id int) (string, error) {
	var name string
	err := config.DB.QueryRow("SELECT name FROM project_category WHERE id=?", id).Scan(&name)

	if err != nil {
		return "", errors.New("Server is unable to execute query to the database")
	}

	return name, nil
}
