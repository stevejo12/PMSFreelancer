package controller

import (
	"PMSFreelancer/config"
	"PMSFreelancer/models"
	"errors"
	"net/http"

	// "github.com.stevejo12/PMSFreelancer/config"
	// "github.com.stevejo12/PMSFreelancer/models"

	"github.com/gin-gonic/gin"
)

func AddUserPortfolio(c *gin.Context) {
	id := c.Param("id")

	// accept the images and store it in the tempfile
	c.Request.ParseMultipartForm(5 * 1024 * 1024)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to read the uploaded file"})
		return
	}
	defer file.Close()

	description := c.Request.FormValue("description")

	if description == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Description should not be empty"})
		return
	}

	url, err := uploadFile(file, header)

	_, err = config.DB.Query("INSERT INTO portfolio(description, link, user_id) VALUES(?,?,?)", description, url, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to read the uploaded file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully added user portfolio"})
}

func EditUserPortfolio(c *gin.Context) {
	id := c.Param("id")

	type structData struct {
		Description string
	}

	var data structData

	err = c.Bind(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	_, err = config.DB.Exec("UPDATE portfolio SET description=? WHERE id=?", data.Description, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to delete the data in the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully edited user portfolio"})
}

func DeleteUserPortfolio(c *gin.Context) {
	id := c.Param("id")

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

func allUserPortfolio(id string) ([]models.PortfolioDatabase, error) {
	data, err := config.DB.Query("SELECT * FROM portfolio WHERE user_id=?", id)

	if err != nil {
		return []models.PortfolioDatabase{}, errors.New(err.Error())
	}

	var returnValue []models.PortfolioDatabase

	for data.Next() {
		var dbData models.PortfolioDatabase
		if err := data.Scan(&dbData.ID, &dbData.Description, &dbData.Link, &dbData.OwnerID); err != nil {
			return []models.PortfolioDatabase{}, errors.New("Something is wrong with the database data")
		}

		returnValue = append(returnValue, dbData)
	}

	return returnValue, nil
}

// maybe useful later or delete
// func GetUserPortfolio(c *gin.Context) {
// 	id := c.Param("id")

// 	allData, err := allUserPortfolio(id)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"code":    http.StatusInternalServerError,
// 			"message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"code":    http.StatusOK,
// 		"message": "Successfully added user portfolio",
// 		"data":    allData})
// }
