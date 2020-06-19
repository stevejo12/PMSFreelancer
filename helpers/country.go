package helpers

import (
	"errors"

	"github.com/stevejo12/PMSFreelancer/config"
	// "PMSFreelancer/config"
)

// CountryName => get the country name based on the id
func getCountryName(id int) (string, error) {
	var countryName string

	err := config.DB.QueryRow("SELECT country_name FROM app_countries WHERE id=?", id).Scan(&countryName)

	if err != nil {
		return "", errors.New("Server is unable to execute query to the database")
	}

	return countryName, nil
}
