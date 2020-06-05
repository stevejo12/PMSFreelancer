package main

import (
	"fmt"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/controller"
	_ "github.com/stevejo12/PMSFreelancer/docs"

	// "PMSFreelancer/config"
	// "PMSFreelancer/controller"
	// _ "PMSFreelancer/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	config.ConnectToCloudinary()
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

	v1 := r.Group("/v1")
	{
		v1.GET("/", handleHome)
		v1.POST("/register", controller.RegisterUserWithPassword)
		v1.POST("/login", controller.LoginUserWithPassword)
		v1.POST("/logout", controller.HandleLogout)
		v1.POST("/createBoardTrello/:id", controller.AuthenticationToken, controller.CreateNewBoard)
		v1.GET("/googleLogin", controller.HandleLoginGoogle)
		v1.GET("/signin-callback", controller.HandleCallbackGoogle)
		v1.PUT("/change-password", controller.AuthenticationToken, controller.ChangeUserPassword)
		v1.GET("/allSkills", controller.AuthenticationToken, controller.GetAllSkills)
		v1.POST("/updateSkills/:id", controller.AuthenticationToken, controller.UpdateUserSkills)
		v1.POST("/resetPassword", controller.ResetPassword)
		v1.POST("/uploadImage", controller.UploadImage)
		v1.POST("/uploadFile", controller.UploadFile)
		v1.GET("/searchProject", controller.SearchProject)
		v1.GET("/userEducation/:id", controller.UserEducation)
		v1.POST("/addEducation/:id", controller.AuthenticationToken, controller.AddEducation)
		v1.GET("/userExperience/:id", controller.UserExperience)
		v1.POST("/addExperience/:id", controller.AddExperience)
		v1.POST("/addProject/:id", controller.AddProject)
		v1.GET("/userProjects/:id", controller.GetAllUserProjects)
		v1.GET("/projectDetail/:id", controller.ProjectDetail)
		v1.POST("/submitProjectInterest/:id", controller.SubmitProjectInterest)
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

	const html = `<html><body><a href="/v1/googleLogin"> Google Log In</a>
	<form enctype="multipart/form-data" action="http://localhost:8080/v1/uploadImage" method="post">
    <input type="file" name="myFile" />
		<input type="submit" value="upload" />
	</form>
	<form enctype="multipart/form-data" action="http://localhost:8080/v1/uploadFile" method="post">
    <input type="file" name="myFile" />
    <input type="submit" value="upload" />
  </form></body></html>`
	c.Writer.Write([]byte(html))
}

// Cors => allow access to non origin
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
