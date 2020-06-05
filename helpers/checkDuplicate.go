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
		return errors.New("Server unable to create your account")
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
		return errors.New("Server unable to create your account")
	}
}
