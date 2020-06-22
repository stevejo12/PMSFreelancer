package controller

import (
	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"

	"github.com/gin-gonic/gin"
)

// AddUserPortfolio => Adding User Portfolio
// AddUserPortfolio godoc
// @Summary Adding User Portfolio
// @Accept  json
// @Tags Portfolio
// @Param token header string true "Token Header"
// @Param Parameters body models.PortfolioRequestParameter true "New Portfolio Description"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /addPortfolio [post]
func AddUserPortfolio(c *gin.Context) {
	id := idToken

	param := models.PortfolioRequestParameter{}

	err := c.BindJSON(&param)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	if (param.StartYear > param.EndYear) || (!helpers.IsYearConsistFourNumber(param.StartYear, param.EndYear)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Start or/and End year is invalid"})
		return
	}

	_, err = config.DB.Exec("INSERT INTO portfolio(title, description, link, user_id, start_year, end_year) VALUES(?,?,?,?,?,?)", param.Title, param.Description, param.Link, id, param.StartYear, param.EndYear)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to execute query to the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully added user portfolio"})
}

// EditUserPortfolio => Updating User Portfolio
// EditUserPortfolio godoc
// @Summary Updating User Portfolio
// @Accept  json
// @Tags Portfolio
// @Param token header string true "Token Header"
// @Param id path int64 true "Project ID"
// @Param Description body models.PortfolioRequestParameter true "Update Portfolio Description"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /editPortfolio/{id} [put]
func EditUserPortfolio(c *gin.Context) {
	id := c.Param("id")

	data := models.PortfolioRequestParameter{}

	err = c.BindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	query := "UPDATE portfolio SET description=\"" + data.Description + "\""
	query = query + ", title=\"" + data.Title + "\""
	query = query + ", link=\"" + data.Link + "\""
	query = query + ", start_year=" + strconv.Itoa(data.StartYear)
	query = query + ", end_year=" + strconv.Itoa(data.EndYear)
	query = query + " WHERE id=" + id

	_, err = config.DB.Exec(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to update the data in the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully update user portfolio"})
}

// DeleteUserPortfolio => Deleting User Portfolio
// DeleteUserPortfolio godoc
// @Summary Deleting User Portfolio
// @Accept  json
// @Tags Portfolio
// @Param token header string true "Token Header"
// @Param id path int64 true "Project ID"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /deletePortfolio/{id} [delete]
func DeleteUserPortfolio(c *gin.Context) {
	id := c.Param("id")

	// check if the portfolio id exist
	dataID, err := config.DB.Query("SELECT * FROM portfolio WHERE id=?", id)

	if !dataID.Next() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "The project ID doesn't exist in the database"})
		return
	}

	_, err = config.DB.Exec("DELETE FROM portfolio WHERE id=?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to delete the data in the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully deleted user portfolio"})
}

func allUserPortfolio(id string) ([]models.PortfolioReturnParameter, error) {
	data, err := config.DB.Query("SELECT id, title, description, link, start_year, end_year  FROM portfolio WHERE user_id=? ORDER BY start_year DESC, end_year DESC, id DESC", id)

	if err != nil {
		return []models.PortfolioReturnParameter{}, errors.New(err.Error())
	}

	returnValue := []models.PortfolioReturnParameter{}

	for data.Next() {
		var dbData models.PortfolioReturnParameter
		if err := data.Scan(&dbData.ID, &dbData.Title, &dbData.Description, &dbData.Link, &dbData.StartYear, &dbData.EndYear); err != nil {
			return []models.PortfolioReturnParameter{}, errors.New("Something is wrong with the database data")
		}

		returnValue = append(returnValue, dbData)
	}

	return returnValue, nil
}
