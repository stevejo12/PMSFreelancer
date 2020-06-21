package helpers

import (
	"errors"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
)

// GetCountryInformation => get the country name based on the id
func GetCountryInformation(id int) (models.CountryDataProfile, error) {
	countryInfo := models.CountryDataProfile{}

	err := config.DB.QueryRow("SELECT id, country_name FROM app_countries WHERE id=?", id).Scan(&countryInfo.ID, &countryInfo.CountryName)

	if err != nil {
		return models.CountryDataProfile{}, errors.New("Server is unable to execute query to the database")
	}

	return countryInfo, nil
}
