package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	// "github.com/stevejo12/PMSFreelancer/config"
	// "github.com/stevejo12/PMSFreelancer/controller"

	// for development
	"PMSFreelancer/config"
	"PMSFreelancer/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// _ "github.com/stevejo12/PMSFreelancer/docs"

	_ "PMSFreelancer/docs"
)

var err error

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

var jwtKey = []byte("key_spirits")

type loginInfo struct {
	username string
	password string
}

func init() {
	config.ConnectToDB()
	config.LoadConfig()
}

// @title Swagger API
// @version 1.0
// @description Swagger API for Golang Project PMS + Freelancer.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support

// @BasePath /v1
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(Cors())

	r.GET("/", handleHome)

	v1 := r.Group("/v1")
	{
		// v1.GET("/", handleHome)
		v1.POST("/register", controller.RegisterUserWithPassword)
		v1.POST("/login", controller.LoginUserWithPassword)
		v1.POST("/logout", controller.HandleLogout)
		v1.POST("/createBoardTrello", controller.AuthenticationToken, controller.CreateNewBoard)
		v1.GET("/googleLogin", handleLoginGoogle)
		v1.GET("/signin-callback", handleCallback)
		v1.PUT("/change-password", controller.AuthenticationToken, controller.ChangeUserPassword)
		v1.POST("/resetPassword", controller.ResetPassword)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run(":8080")

	// ini untuk swagger
	// reference : https://golangexample.com/automatically-generate-restful-api-documentation-with-swagger-2-0-for-go/
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err != nil {
		log.Fatal(err)
	}
}

// handleHome, handleLoginGoogle, handleCallback might change later
// this functions work for google sign in will be getting the google_id for the email
func handleHome(c *gin.Context) {
	if c.Request.URL.Path != "/v1/" {
		fmt.Println("error wrong path")
		return
	}

	const html = `<html><body><a href="/v1/googleLogin"> Google Log In</a></body></html>`
	c.Writer.Write([]byte(html))
}

func handleLoginGoogle(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleCallback(c *gin.Context) {
	if c.Request.FormValue("state") != randomState {
		fmt.Println("State is not valid")
		c.Redirect(http.StatusTemporaryRedirect, "/v1")
		c.Abort()
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, c.Request.FormValue("code"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		c.Abort()
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		c.Abort()
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		c.Abort()
	}

	fmt.Fprintf(c.Writer, "Response: %s", content)
}

// Cors => allow access to non origin
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
