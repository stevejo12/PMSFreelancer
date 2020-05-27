package controller

import (
	"PMSFreelancer/config"
	"PMSFreelancer/models"
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
