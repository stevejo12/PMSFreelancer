package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
)

func userExperience(id string) ([]models.ExperienceReturnValue, error) {
	resp, err := config.DB.Query("SELECT * FROM experience WHERE user_id=?", id)

	if err != nil {
		return []models.ExperienceReturnValue{}, errors.New("Server unable to execute query to database")
	}

	var allData []models.ExperienceReturnValue

	for resp.Next() {
		var databaseData models.ExperienceTableResponse
		if err := resp.Scan(&databaseData.ID, &databaseData.Place, &databaseData.Position, &databaseData.StartYear, &databaseData.EndYear, &databaseData.UserID); err != nil {
			return []models.ExperienceReturnValue{}, errors.New("Something is wrong with the database data")
		}

		var returnValue models.ExperienceReturnValue

		returnValue.ID = databaseData.ID
		returnValue.Position = databaseData.Position
		returnValue.Place = databaseData.Place
		returnValue.StartYear = databaseData.StartYear
		returnValue.EndYear = databaseData.EndYear

		allData = append(allData, returnValue)
	}

	if resp.Err() != nil {
		return []models.ExperienceReturnValue{}, errors.New("Something is wrong with the data retrieved")
	}

	return allData, nil
}

func GetOnlyUserExperience(c *gin.Context) {
	id := c.Param("id")

	allUserExperience, err := userExperience(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    []models.ExperienceReturnValue{}})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All User Experience data have been retrieved",
		"data":    allUserExperience})
}

func AddExperience(c *gin.Context) {
	id := c.Param("id")

	// sample data
	// experience: [
	// 	{
	// 		Place
	//		Position
	// 		StartYear
	// 		EndYear
	// 	},
	// 	{
	// 		Place
	//		Position
	// 		StartYear
	// 		EndYear
	// 	}
	// ]

	var data models.AddExperienceParameter

	err = c.Bind(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	experience := data.Experience

	query := "INSERT INTO experience(place, position, starting_year, ending_year, user_id) VALUES"
	for i := 0; i < len(experience); i++ {
		place := experience[i].Place
		position := experience[i].Position
		startingYear := experience[i].StartYear
		endYear := experience[i].EndYear
		query = query + "(\"" + place + "\", \"" + position + "\"," + strconv.Itoa(startingYear) + ", " + strconv.Itoa(endYear) + ", " + id + "),"
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
		"message": "Successfully Added Experience"})
}
