package helpers

import (
	"github.com/stevejo12/PMSFreelancer/config"
	// "PMSFreelancer/config"
	"errors"
)

func CountryList(country_id int) error {
	rows, err := config.DB.Query("SELECT id from app_countries")

	if err != nil {
		return err
	}

	var allCountries []interface{}
	var countryid int

	for rows.Next() {
		err := rows.Scan(&countryid)

		if err != nil {
			return err
		}

		allCountries = append(allCountries, countryid)
	}

	exist := Contains(allCountries, country_id)

	if !exist {
		return errors.New("not exist")
	}

	return nil
}
