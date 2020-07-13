package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"google.golang.org/api/oauth2/v2"
)

var httpClient = &http.Client{}

func googleLoginVerification(c *gin.Context) {
	oauth2Service, err := oauth2.New(httpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Google token is not valid"})
	}

	defer httpClient.CloseIdleConnections()

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Login Successful",
		"data":    tokenInfo})
}
