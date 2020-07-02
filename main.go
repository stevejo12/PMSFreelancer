package main

import (
	"fmt"
	"log"

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

func init() {
	config.ConnectToDB()
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
		v1.POST("/checkEmail", controller.CheckEmail)
		v1.POST("/register", controller.RegisterUserWithPassword)
		v1.GET("/googleRegister", controller.HandleRegisterGoogle)
		v1.GET("/registerCallback", controller.HandleCallbackRegisterGoogle)
		v1.POST("/registerUserUsingGoogle", controller.RegisterUserWithGoogle)
		v1.POST("/login", controller.LoginUserWithPassword)
		v1.POST("/logout", controller.AuthenticationToken, controller.HandleLogout)
		v1.GET("/allCountries", controller.GetAllCountries)
		v1.GET("/userProfile/:id", controller.AuthenticationToken, controller.GetUserProfileByID)
		v1.GET("/userProfile", controller.AuthenticationToken, controller.GetUserProfile)
		v1.POST("/userReview", controller.AuthenticationToken, controller.AddUserReview)
		v1.PUT("/updateProfile", controller.AuthenticationToken, controller.UpdateUserProfile)
		v1.POST("/addPortfolio", controller.AuthenticationToken, controller.AddUserPortfolio)
		v1.DELETE("/deletePortfolio/:id", controller.AuthenticationToken, controller.DeleteUserPortfolio)
		v1.PUT("/editPortfolio/:id", controller.AuthenticationToken, controller.EditUserPortfolio)
		v1.POST("/createBoardTrello", controller.AuthenticationToken, controller.CreateNewBoard)
		v1.GET("/googleLogin", controller.HandleLoginGoogle)
		v1.GET("/signin-callback", controller.HandleCallbackLoginGoogle)
		v1.PUT("/change-password", controller.AuthenticationToken, controller.ChangeUserPassword)
		v1.GET("/allSkills", controller.GetAllSkills)
		v1.PUT("/editSkills", controller.AuthenticationToken, controller.UpdateUserSkills)
		v1.POST("/resetPassword", controller.ResetPassword)
		v1.POST("/updateNewPassword", controller.UpdateNewPassword)
		v1.POST("/uploadPicture", controller.AuthenticationToken, controller.UploadPicture)
		v1.POST("/uploadAttachment", controller.AuthenticationToken, controller.UploadAttachment)
		v1.GET("/searchProject", controller.SearchProject)
		v1.POST("/addEducation", controller.AuthenticationToken, controller.AddEducation)
		v1.PUT("/editEducation/:id", controller.AuthenticationToken, controller.EditEducation)
		v1.DELETE("/deleteEducation/:id", controller.AuthenticationToken, controller.DeleteEducation)
		v1.POST("/addExperience", controller.AuthenticationToken, controller.AddExperience)
		v1.PUT("/editExperience/:id", controller.AuthenticationToken, controller.EditExperience)
		v1.DELETE("/deleteExperience/:id", controller.AuthenticationToken, controller.DeleteExperience)
		v1.POST("/addProject", controller.AuthenticationToken, controller.AddProject)
		v1.PUT("/editProject/:id", controller.AuthenticationToken, controller.EditProject)
		v1.DELETE("/deleteProject/:id", controller.AuthenticationToken, controller.DeleteProject)
		v1.GET("/userProjects", controller.AuthenticationToken, controller.GetAllUserProjects)
		v1.GET("/projectDetail/:id", controller.ProjectDetail)
		v1.POST("/submitProjectInterest/:id", controller.AuthenticationToken, controller.SubmitProjectInterest)
		v1.POST("/acceptProjectInterest/:id", controller.AuthenticationToken, controller.AcceptProjectInterest)
		v1.POST("/submitProjectForReview/:id", controller.AuthenticationToken, controller.ReviewProject)
		v1.POST("/rejectReviewProject/:id", controller.AuthenticationToken, controller.RejectReviewProject)
		v1.POST("/completeProject/:id", controller.AuthenticationToken, controller.CompleteProject)
		v1.GET("/allProjectCategory", controller.AuthenticationToken, controller.GetAllProjectCategory)
		v1.POST("/depositMoney", controller.AuthenticationToken, controller.DepositMoney)
		v1.POST("/checkMutationMoota", controller.GetTransactionMutation)
		v1.POST("/submitWithdrawRequest", controller.AuthenticationToken, controller.WithdrawMoney)
		v1.GET("/allWithdrawRequests", controller.AuthenticationToken, controller.GetAllWithdrawRequest)
		v1.GET("/userWithdrawRequest", controller.AuthenticationToken, controller.GetUserWithdrawRequest)
		v1.DELETE("/deleteUserWithdrawRequest/:id", controller.AuthenticationToken, controller.DeleteWithdrawRequest)
		v1.PUT("/completeWithdrawRequest/:id", controller.AuthenticationToken, controller.CompleteWithdrawRequest)
		// v1.GET("/testAPI", controller.TestFunction)
		// v1.GET("/mootaProfile", controller.GetMootaProfile)
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

	const html = `<html><body><a href="/v1/googleLogin"> Google Log In</a><a href="/v1/googleRegister"> Google Register</a>
	<form enctype="multipart/form-data" action="http://159.89.202.223:8080/v1/uploadImage" method="post">
    <input type="file" name="file" />
		<input type="submit" value="upload" />
	</form>
	<form enctype="multipart/form-data" action="http://159.89.202.223:8080/v1/uploadAttachment" method="post">
		<input type="file" name="file" />
    <input type="submit" value="upload" />
  </form></body></html>`
	c.Writer.Write([]byte(html))
}

// Cors => allow access to non origin
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
