package helpers

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// ConvertDate => convert MySQL format something like "20 June 2020"
func ConvertDate(date string) (string, error) {
	// expected date format YYYY-MM-DDT00:00:00Z
	// therefore we want to get the YYYY-MM-DD only
	splitDateTime := strings.Split(date, "T")

	if len(splitDateTime) == 2 {
		format := "2006-01-02"
		// get the YYYY-MM-DD
		getYMD := splitDateTime[0]

		splittedDate := strings.Split(getYMD, "-")

		if len(splittedDate) == 3 {
			dateTime, err := time.Parse(format, getYMD)

			if err != nil {
				return date, errors.New("Parsing error")
			}

			// Format data 20 June 2020
			newDateFormat := strconv.Itoa(dateTime.Day()) + " " + dateTime.Month().String() + " " + strconv.Itoa(dateTime.Year())

			return newDateFormat, nil
		}

		return date, errors.New("Database date format is not as expected")
	}

	return date, errors.New("Database date format is not as expected")
}
