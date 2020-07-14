package helpers

import (
	// "PMSFreelancer/models"
	"github.com/stevejo12/PMSFreelancer/models"
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

// GetCategoryRaw => get the id and name of the category
func GetCategoryRaw(id int) (models.CategoryRaw, error) {
	var categoryRaw models.CategoryRaw
	err := config.DB.QueryRow("SELECT id, name FROM project_category WHERE id=?", id).Scan(&categoryRaw.ID, &categoryRaw.Name)

	if err != nil {
		return models.CategoryRaw{}, errors.New("Server is unable to execute query to the database")
	}

	return categoryRaw, nil
}
