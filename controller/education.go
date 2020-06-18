package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
)

func userEducation(id string) ([]models.EducationReturnValue, error) {
	returnData := []models.EducationReturnValue{}

	resp, err := config.DB.Query("SELECT * FROM education WHERE user_id=?", id)

	if err != nil {
		return returnData, errors.New("Server unable to execute query to database")
	}

	for resp.Next() {
		var databaseData models.EducationTableResponse
		if err := resp.Scan(&databaseData.ID, &databaseData.Name, &databaseData.StartYear, &databaseData.EndYear, &databaseData.UserID); err != nil {
			return []models.EducationReturnValue{}, errors.New("Something is wrong with the database data")
		}

		var returnValue models.EducationReturnValue

		returnValue.ID = databaseData.ID
		returnValue.Name = databaseData.Name
		returnValue.StartYear = databaseData.StartYear
		returnValue.EndYear = databaseData.EndYear

		returnData = append(returnData, returnValue)
	}

	if resp.Err() != nil {
		return []models.EducationReturnValue{}, errors.New("Something is wrong with the data retrieved")
	}

	return returnData, nil
}

// GetOnlyUserEducation => Get Detail View for the User Education
// GetOnlyUserEducation godoc
// @Summary User Education
// @Produce json
// @Accept  json
// @Tags Education
// @Param token header string true "Token Header"
// @Success 200 {object} models.ResponseOKGetUserEducation
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /userEducation [get]
func GetOnlyUserEducation(c *gin.Context) {
	id := idToken

	allUserEducation, err := userEducation(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    []models.EducationReturnValue{}})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All User Education data have been retrieved",
		"data":    allUserEducation})
}

// AddEducation => Add User Education
// AddEducation godoc
// @Summary Add User Education
// @Produce json
// @Accept  json
// @Tags Education
// @Param token header string true "Token Header"
// @Param Data body models.AddEducationParameter true "Data Format to add education"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /addEducation [post]
func AddEducation(c *gin.Context) {
	id := idToken

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

	err = c.BindJSON(&data)

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully Added Education"})
}
