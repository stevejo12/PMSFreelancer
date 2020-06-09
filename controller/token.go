package controller

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stevejo12/PMSFreelancer/models"
)

func generateToken(email string) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.TokenClaims{
		Username: email,
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
