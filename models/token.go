package models

import "github.com/dgrijalva/jwt-go"

// TokenClaims => store token to specific username
type TokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
