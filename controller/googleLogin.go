package controller

import (
	"PMSFreelancer/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/v1/signin-callback",
		ClientID:     "776281301027-aincdrlljhjdmu39lfq2aunqeofn1hi8.apps.googleusercontent.com",
		ClientSecret: "5q_niwCvO1dAFEzT2QkcQkok",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	// TO DO: randomize it
	randomState = "random"
)

func HandleLoginGoogle(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleCallbackGoogle(c *gin.Context) {
	if c.Request.FormValue("state") != randomState {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "State value is invalid",
		})
		return
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, c.Request.FormValue("code"))

	if err != nil {
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

	fmt.Println(account.ID)
	fmt.Println(account.Email)
	fmt.Println(account.Verified_email)
	fmt.Println(account.Picture)

	fmt.Fprintf(c.Writer, "Response: %s", content)
}
