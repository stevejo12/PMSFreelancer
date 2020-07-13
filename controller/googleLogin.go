package controller

import (
	"database/sql"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"

	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfigLogin = &oauth2.Config{
		RedirectURL:  "http://api.spirits.id:8080/v1/signin-callback",
		ClientID:     "149260447643-kpqpukphmhj907876qg6q1rmhhit7831.apps.googleusercontent.com",
		ClientSecret: "0OWYl4x5-74O2Mx5AiZgndC_",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	// TO DO: randomize it
	randomState = "random"
)

func HandleLoginGoogle(c *gin.Context) {
	url := googleOauthConfigLogin.AuthCodeURL(randomState)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleCallbackLoginGoogle(c *gin.Context) {
	if c.Request.FormValue("state") != randomState {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "State value is invalid",
		})
		return
	}

	token, err := googleOauthConfigLogin.Exchange(oauth2.NoContext, c.Request.FormValue("code"))

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "There is a problem exchanging token with google",
		})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Unable to access google API",
		})
		c.Abort()
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Error in reading the response from google",
		})
		c.Abort()
	}

	var account models.GoogleResponse

	err = json.Unmarshal(content, &account)

	if err != nil {
		fmt.Println(err)
		return
	}

	// find the email
	var email string
	var id string
	err = config.DB.QueryRow("SELECT id, email FROM login WHERE email=?", account.Email).Scan(&id, &email)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Email is not registered in our database"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to get information from database"})
		return
	}

	cookieToken, expirationTime, err := generateToken(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   cookieToken,
	// 	Expires: expirationTime,
	// })

	tokenInfo := models.TokenResponse{}
	tokenInfo.Token = cookieToken
	tokenInfo.Expire = expirationTime

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Login Successful",
		"data":    tokenInfo})
}
