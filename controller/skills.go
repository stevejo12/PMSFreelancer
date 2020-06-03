package controller

import (
	"errors"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"

	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
	// "PMSFreelancer/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllSkills => Retrieved all the possible skill options available for users
func GetAllSkills(c *gin.Context) {
	data, err := config.DB.Query("SELECT * FROM skills")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Server unable to execute query to database"})
		return
	}

	var allData []models.UserSkills

	for data.Next() {
		// Scan one customer record
		var skills models.UserSkills
		if err := data.Scan(&skills.ID, &skills.Name, &skills.Created_at, &skills.Updated_at); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusText(http.StatusInternalServerError),
				"message": "Something is wrong with the database data"})
			return
		}
		allData = append(allData, skills)
	}
	if data.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Something is wrong with the data retrieved"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Server unable to execute query"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusText(http.StatusOK),
		"message": "All Skills data have been successfully retrieved",
		"data":    allData})
}

func getSkillNames(param string) ([]string, error) {
	var result []string
	initialQuery, err := helpers.SettingInQueryWithID("skills", param)

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
		// var name string
		// if err := data.Scan(&name); err != nil {
		// 	return []string{}, errors.New("Something is wrong with the database data")
		// }
		// fmt.Println(name)
		// result = append(result, name)
	}
	if data.Err() != nil {
		return []string{}, errors.New("Something is wrong with the data retrieved")
	}

	return result, nil
}
