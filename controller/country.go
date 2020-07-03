package controller

import (
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"

	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllCountries godoc
// @Summary Getting all list of countries
// @Produce json
// @Tags Location
// @Success 200 {object} models.ResponseOKGetAllCountries
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /allCountries [get]
func GetAllCountries(c *gin.Context) {
	rows, err := config.DB.Query("SELECT * from app_countries")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with query to get the country list"})
		return
	}

	var allCountries []models.CountryData

	for rows.Next() {
		var country models.CountryData
		if err := rows.Scan(&country.ID, &country.CountryCode, &country.CountryName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		allCountries = append(allCountries, country)
	}

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Country List has been retrieved",
		"data":    allCountries})
}
