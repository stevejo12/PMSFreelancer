package controller

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stevejo12/PMSFreelancer/models"
	// "PMSFreelancer/models"
)

func generateToken(userID string) (string, time.Time, error) {
	expirationTime := time.Now().Add(30 * time.Minute)

	var isAdmin bool
	if userID == "43" {
		isAdmin = true
	} else {
		isAdmin = false
	}

	claims := &models.TokenClaims{
		Username: userID,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(jwtKey)

	if err != nil {
		return "", time.Time{}, errors.New("Server is unable to generate token")
	}

	return tokenString, expirationTime, nil
}
