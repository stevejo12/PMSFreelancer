package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenClaims => store token to specific username
type TokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type TokenResponse struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}
