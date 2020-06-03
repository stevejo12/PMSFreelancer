package controller

import (
	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"
	"database/sql"
	"errors"
	"net/http"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"

	"github.com/gin-gonic/gin"
)

// SearchProject => Search project in SPIRITS
func SearchProject(c *gin.Context) {
	// initialize variables
	// page is page number in pagination
	// size is the number of result per page
	// var page = 0
	page := 0
	size := 10

	// record #1 is number 0
	var startingRecordNumber = page * size
	var endingRecordNumber = startingRecordNumber + (size - 1)

	result, err := config.DB.Query("SELECT * FROM project ORDER BY ID ASC LIMIT ?,?", startingRecordNumber, endingRecordNumber)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with query to get the project list"})
		return
	}

	var allData []models.ProjectSearchResponse
	for result.Next() {
		var project models.ProjectResult
		var data models.ProjectSearchResponse
		var s sql.NullString
		if err := result.Scan(&project.ID, &project.Title, &project.Skills, &project.Price, &s); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		if s.Valid {
			data.Attachment, err = getProjectLinks(s.String)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Something is wrong with the database data"})
				return
			}
		} else {
			// if the data in the database is null attachment validity (valid) will be false
			data.Attachment = []string{}
		}

		data.ID = project.ID
		data.Title = project.Title
		data.Skills, err = getSkillNames(project.Skills)
		data.Price = project.Price
		allData = append(allData, data)
	}

	if result.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with the data retrieved"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All Project data have been retrieved",
		"data":    allData})
}

func getProjectLinks(param string) ([]string, error) {
	var result []string
	query, err := helpers.SettingInQueryWithID("project_links", param)

	if err != nil {
		return nil, err
	}

	data, err := config.DB.Query(query)

	if err != nil {
		return nil, errors.New("Server unable to execute query to database")
	}

	for data.Next() {
		// Scan one customer record
		var link models.ProjectLinksResponse
		if err := data.Scan(&link.ID, &link.Project_link); err != nil {
			return []string{}, errors.New("Something is wrong with the database data")
		}
		result = append(result, link.Project_link)
	}
	if data.Err() != nil {
		return []string{}, errors.New("Something is wrong with the data retrieved")
	}

	return result, nil
}
