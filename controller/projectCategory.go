package controller

import (
	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
)

// GetAllProjectCategory => Get the list of possible project category
// GetAllProjectCategory godoc
// @Summary Get All Project Category List
// @Produce json
// @Tags Project
// @Param token header string true "Token Header"
// @Success 200 {object} models.ResponseOKProjectCategory
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /allProjectCategory [get]
func GetAllProjectCategory(c *gin.Context) {
	allData := []models.ProjectCategoryData{}
	result, err := config.DB.Query("SELECT * FROM project_category")
	defer result.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to execute query to the database"})
		return
	}

	for result.Next() {
		data := models.ProjectCategoryData{}
		if err := result.Scan(&data.ID, &data.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		allData = append(allData, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Fetching Project Category List successful",
		"data":    allData})
}
