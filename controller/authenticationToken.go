package controller

import (
	// "PMSFreelancer/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/models"
)

var idToken string
var isAdmin bool

// AuthenticationToken => authentication to get a token for login into SPIRITS
func AuthenticationToken(c *gin.Context) {
	cookieStrSwagger := c.Request.Header.Get("token")
	cookie, err := c.Request.Cookie("token")

	// to accomodate swagger cookie => get the token from the header
	// gin swagger is OPenAPI swagger 2.0
	// can not support cookie
	// no library yet for 3.0 version
	if cookieStrSwagger != "" && cookie == nil {
		claimsSwagger := &models.TokenClaims{}

		tknSwagger, errSwagger := jwt.ParseWithClaims(cookieStrSwagger, claimsSwagger, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if errSwagger != nil {
			if !tknSwagger.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Token is no longer invalid",
				})
				c.Abort()
				return
			}

			if errSwagger == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Token is invalid",
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		idToken = claimsSwagger.Username
		isAdmin = claimsSwagger.IsAdmin
		err = nil
	}

	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Token is not present",
			})
			c.Abort()
			return
		}

		// For any other type of error, return a bad request status
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	// for regular cookie
	if cookie != nil {
		// Get the JWT string from the cookie
		tknStr := cookie.Value

		// Initialize a new instance of `Claims`
		claims := &models.TokenClaims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if !tkn.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Token is no longer invalid",
				})
				c.Abort()
				return
			}

			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Token is invalid",
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		idToken = claims.Username
	}
}
