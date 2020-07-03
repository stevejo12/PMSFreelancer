package controller

import (
	"errors"
	"net/http"
	"strconv"

	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
	// "PMSFreelancer/helpers"

	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"
)

func userExperience(id string) ([]models.ExperienceReturnValue, error) {
	resp, err := config.DB.Query("SELECT * FROM experience WHERE user_id=? ORDER BY starting_year DESC, ending_year DESC, id DESC", id)

	defer resp.Close()

	if err != nil {
		return []models.ExperienceReturnValue{}, errors.New("Server unable to execute query to database")
	}

	allData := []models.ExperienceReturnValue{}

	for resp.Next() {
		var databaseData models.ExperienceTableResponse
		if err := resp.Scan(&databaseData.ID, &databaseData.Description, &databaseData.Place, &databaseData.Position, &databaseData.StartYear, &databaseData.EndYear, &databaseData.UserID); err != nil {
			return []models.ExperienceReturnValue{}, errors.New("Something is wrong with the database data")
		}

		var returnValue models.ExperienceReturnValue

		returnValue.ID = databaseData.ID
		returnValue.Position = databaseData.Position
		returnValue.Place = databaseData.Place
		returnValue.StartYear = databaseData.StartYear
		returnValue.EndYear = databaseData.EndYear
		returnValue.Description = databaseData.Description

		allData = append(allData, returnValue)
	}

	if resp.Err() != nil {
		return []models.ExperienceReturnValue{}, errors.New("Something is wrong with the data retrieved")
	}

	return allData, nil
}

// AddExperience => Add User Experience
// AddExperience godoc
// @Summary Adding User Experience
// @Produce json
// @Accept  json
// @Tags Experience
// @Param token header string true "Token Header"
// @Param Data body models.ExperienceParameters true "Data Format to add experience"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /addExperience [post]
func AddExperience(c *gin.Context) {
	id := idToken

	var data models.ExperienceParameters

	err = c.BindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	if (data.StartYear > data.EndYear) || (!helpers.IsYearConsistFourNumber(data.StartYear, data.EndYear)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Start or/and End year is invalid"})
		return
	}

	query := "INSERT INTO experience(place, position, starting_year, ending_year, user_id, description) VALUES"
	query = query + "(\"" + data.Place + "\", \"" + data.Position + "\"," + strconv.Itoa(data.StartYear) + ", " + strconv.Itoa(data.EndYear) + ", " + id + ", \"" + data.Description + "\")"

	_, err = config.DB.Exec(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully Added Experience"})
}

// EditExperience => Updating User Experience
// EditExperience godoc
// @Summary Updating User Experience
// @Accept  json
// @Tags Experience
// @Param token header string true "Token Header"
// @Param id path int64 true "Experience ID"
// @Param Description body models.ExperienceParameters true "New Experience Description"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /editExperience/{id} [put]
func EditExperience(c *gin.Context) {
	id := c.Param("id")

	data := models.ExperienceParameters{}

	err = c.BindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	// setting up query so it is not so long per line
	query := "UPDATE experience SET place=\"" + data.Place + "\""
	query = query + ", position=\"" + data.Position + "\""
	query = query + ", starting_year=\"" + strconv.Itoa(data.StartYear) + "\""
	query = query + ", ending_year=\"" + strconv.Itoa(data.EndYear) + "\""
	query = query + ", description=\"" + data.Description + "\""
	query = query + " WHERE id=" + id

	_, err = config.DB.Exec(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to execute query to the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully edited user Experience"})
}

// DeleteExperience => Deleting User Experience
// DeleteExperience godoc
// @Summary Deleting User Experience
// @Accept  json
// @Tags Experience
// @Param token header string true "Token Header"
// @Param id path int64 true "Project ID"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /deleteExperience/{id} [delete]
func DeleteExperience(c *gin.Context) {
	id := c.Param("id")

	// check if the Experience id exist
	dataID, err := config.DB.Query("SELECT * FROM experience WHERE id=?", id)
	defer dataID.Close()

	if !dataID.Next() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "The Experience ID doesn't exist in the database"})
		return
	}

	_, err = config.DB.Exec("DELETE FROM experience WHERE id=?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to delete the data in the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully deleted user Experience"})
}
