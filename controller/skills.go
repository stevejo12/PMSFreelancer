package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"
	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"
)

func getAllSkills() ([]models.UserSkills, error) {
	data, err := config.DB.Query("SELECT * FROM skills")
	var allData []models.UserSkills

	if err != nil {
		return allData, errors.New("Server unable to execute query to database")
	}

	for data.Next() {
		// Scan one customer record
		var skills models.UserSkills
		if err := data.Scan(&skills.ID, &skills.Name, &skills.Created_at, &skills.Updated_at); err != nil {
			return []models.UserSkills{}, errors.New("Something is wrong with the database data")
		}
		allData = append(allData, skills)
	}
	if data.Err() != nil {
		return []models.UserSkills{}, errors.New("Something is wrong with the data retrieved")
	}

	return allData, nil
}

// GetAllSkills => Get a list of available skills
func GetAllSkills(c *gin.Context) {
	allSkills, err := getAllSkills()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    []models.UserSkills{}})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All Skills data have been retrieved",
		"data":    allSkills})
}

func getSkillNames(param string) ([]string, error) {
	var result []string
	initialQuery, err := helpers.SettingInQueryWithID("skills", param, "*")

	if err != nil {
		return nil, err
	}

	data, err := config.DB.Query(initialQuery)

	if err != nil {
		return result, errors.New("Server unable to execute query to database")
	}

	for data.Next() {
		// Scan one customer record
		var skills models.UserSkills
		if err := data.Scan(&skills.ID, &skills.Name, &skills.Created_at, &skills.Updated_at); err != nil {
			return []string{}, errors.New("Something is wrong with the database data")
		}
		result = append(result, skills.Name)
	}
	if data.Err() != nil {
		return []string{}, errors.New("Something is wrong with the data retrieved")
	}

	return result, nil
}

func userSkills(id string) ([]models.UserSkills, error) {
	query, err := helpers.SettingInQueryWithID("skills", id, "*")

	if err != nil {
		return []models.UserSkills{}, errors.New(err.Error())
	}

	resp, err := config.DB.Query(query)

	if err != nil {
		return []models.UserSkills{}, errors.New(err.Error())
	}

	var allData []models.UserSkills

	for resp.Next() {
		var databaseData models.UserSkills
		if err := resp.Scan(&databaseData.ID, &databaseData.Name, &databaseData.Created_at, &databaseData.Updated_at); err != nil {
			return []models.UserSkills{}, errors.New("Something is wrong with the database data")
		}

		allData = append(allData, databaseData)
	}

	if resp.Err() != nil {
		return []models.UserSkills{}, errors.New("Something is wrong with the data retrieved")
	}

	return allData, nil
}
