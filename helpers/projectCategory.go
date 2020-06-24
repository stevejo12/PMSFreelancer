package helpers

import (
	// "PMSFreelancer/config"
	"errors"

	"github.com/stevejo12/PMSFreelancer/config"
)

// IsThisCategoryIDExist => Check if int value id exist in the category database
func IsThisCategoryIDExist(id int) error {
	var allCategories []interface{}

	rows, err := config.DB.Query("SELECT id from project_category")

	if err != nil {
		return errors.New("Server is unable to execute query to the database")
	}

	var categoryID int
	for rows.Next() {
		err := rows.Scan(&categoryID)

		if err != nil {
			return err
		}

		allCategories = append(allCategories, categoryID)
	}

	exist := Contains(allCategories, id)

	if !exist {
		return errors.New("not exist")
	}

	return nil
}
