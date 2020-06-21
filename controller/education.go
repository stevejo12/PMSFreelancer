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

	resp, err := config.DB.Query("SELECT * FROM education WHERE user_id=? ORDER BY starting_year DESC, ending_year DESC, id DESC", id)

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

// AddEducation => Add User Education
// AddEducation godoc
// @Summary Add User Education
// @Produce json
// @Accept  json
// @Tags Education
// @Param token header string true "Token Header"
// @Param Data body models.EducationParameters true "Data Format to add education"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /addEducation [post]
func AddEducation(c *gin.Context) {
	id := idToken

	var data models.EducationParameters

	err = c.BindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	if data.StartYear > data.EndYear {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Start year should be in the past compared to end year"})
		return
	}

	query := "INSERT INTO education(name, starting_year, ending_year, user_id) VALUES"
	query = query + "(\"" + data.Name + "\", " + strconv.Itoa(data.StartYear) + ", " + strconv.Itoa(data.EndYear) + ", " + id + ")"

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

// EditEducation => Updating User Education
// EditEducation godoc
// @Summary Updating User Education
// @Accept  json
// @Tags Education
// @Param token header string true "Token Header"
// @Param id path int64 true "Education ID"
// @Param Description body models.EducationParameters true "New Education Description"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /editEducation/{id} [put]
func EditEducation(c *gin.Context) {
	id := c.Param("id")

	data := models.EducationParameters{}

	err = c.BindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	query := "UPDATE education SET name=\"" + data.Name + "\", starting_year=" + strconv.Itoa(data.StartYear) + ", ending_year=" + strconv.Itoa(data.EndYear) + " WHERE id=" + id

	_, err = config.DB.Exec(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to execute query to the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully edited user education"})
}

// DeleteEducation => Deleting User Education
// DeleteEducation godoc
// @Summary Deleting User Education
// @Accept  json
// @Tags Education
// @Param token header string true "Token Header"
// @Param id path int64 true "Education ID"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /deleteEducation/{id} [delete]
func DeleteEducation(c *gin.Context) {
	id := c.Param("id")

	// check if the education id exist
	dataID, err := config.DB.Query("SELECT * FROM education WHERE id=?", id)

	if !dataID.Next() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "The education ID doesn't exist in the database"})
		return
	}

	_, err = config.DB.Exec("DELETE FROM education WHERE id=?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to delete the data in the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully deleted user education"})
}
