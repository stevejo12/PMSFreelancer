package helpers

import (
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
	"database/sql"
	"errors"
)

// GetProjectDetailTitle => This is for review section (project data)
func GetProjectDetailTitle(id int) (models.ProjectInformationForReview, error) {
	result := models.ProjectInformationForReview{}
	err := config.DB.QueryRow("SELECT id, title FROM project WHERE id=?", id).Scan(&result.ID, &result.Title)

	if err == sql.ErrNoRows {
		return models.ProjectInformationForReview{}, errors.New("Server can't find the project in the database")
	} else if err != nil {
		return models.ProjectInformationForReview{}, errors.New("Server is unable to execute query")
	}

	return result, nil
}

// GetUserInformationForReview => This is for review section (User data)
func GetUserInformationForReview(id int) (models.UserInformationForReview, error) {
	result := models.UserInformationForReview{}
	err := config.DB.QueryRow("SELECT id, first_name, last_name FROM login WHERE id=?", id).Scan(&result.ID, &result.FirstName, &result.LastName)

	if err == sql.ErrNoRows {
		return models.UserInformationForReview{}, errors.New("Server can't find the user information in the database")
	} else if err != nil {
		return models.UserInformationForReview{}, errors.New("Server is unable to execute query")
	}

	return result, nil
}
