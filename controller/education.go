package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
)

func UserEducation(c *gin.Context) {
	id := c.Param("id")

	resp, err := config.DB.Query("SELECT * FROM education WHERE user_id=?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})
		return
	}

	var allData []models.EducationReturnValue

	for resp.Next() {
		var databaseData models.EducationTableResponse
		if err := resp.Scan(&databaseData.ID, &databaseData.Name, &databaseData.StartYear, &databaseData.EndYear, &databaseData.UserID); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		var returnValue models.EducationReturnValue

		returnValue.ID = databaseData.ID
		returnValue.Name = databaseData.Name
		returnValue.StartYear = databaseData.StartYear
		returnValue.EndYear = databaseData.EndYear

		allData = append(allData, returnValue)
	}

	if resp.Err() != nil {
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

func AddEducation(c *gin.Context) {
	id := c.Param("id")

	// sample data
	// education: [
	// 	{
	// 		Name
	// 		StartYear
	// 		EndYear
	// 	},
	// 	{
	// 		Name
	// 		StartYear
	// 		EndYear
	// 	}
	// ]

	var data models.AddEducationParameter

	err = c.Bind(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	education := data.Education

	query := "INSERT INTO education(name, starting_year, ending_year, user_id) VALUES"
	for i := 0; i < len(education); i++ {
		name := education[i].Name
		startingYear := education[i].StartYear
		endYear := education[i].EndYear
		query = query + "(\"" + name + "\", " + strconv.Itoa(startingYear) + ", " + strconv.Itoa(endYear) + ", " + id + "),"
	}

	if last := len(query) - 1; last >= 0 && query[last] == ',' {
		query = query[:last]
	}

	_, err = config.DB.Exec(query)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully Added Education"})
}
