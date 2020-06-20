package helpers

import (
	"github.com/stevejo12/PMSFreelancer/config"
	// "PMSFreelancer/config"
	"database/sql"
	"errors"
)

func CheckDuplicateEmail(email string) error {
	var databaseEmail string
	row := config.DB.QueryRow("SELECT email FROM login WHERE email=?", email).Scan(&databaseEmail)

	if row == sql.ErrNoRows {
		return nil
	} else if row == nil {
		return errors.New("Email has already registered in the database")
	} else {
		return errors.New(row.Error())
	}
}

func CheckDuplicateUsername(username string) error {
	var databaseUsername string
	row := config.DB.QueryRow("SELECT email FROM login WHERE username=?", username).Scan(&databaseUsername)

	if row == sql.ErrNoRows {
		return nil
	} else if row == nil {
		return errors.New("Username has already registered in the database")
	} else {
		return errors.New(row.Error())
	}
}

func RemoveDuplicateValueArray(arr []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func CheckUpdatedUsername(username string, id string) error {
	var existingID int
	row := config.DB.QueryRow("SELECT id FROM login WHERE username=? && id !=?", username, id).Scan(&existingID)

	if row == sql.ErrNoRows {
		return nil
	} else if row == nil {
		return errors.New("Username has already been taken")
	} else {
		return errors.New(row.Error())
	}
}
