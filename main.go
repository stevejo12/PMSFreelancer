package main

import (
	"database/sql"
	"log"

	"github.com/stevejo12/PMSFreelancer/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/stevejo12/PMSFreelancer/docs"
)

// DB => Store DB Connection upon initialization
var DB *sql.DB
var err error

type loginInfo struct {
	username string
	password string
}

func init() {
	controller.ConnectToDB()
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

	// url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")

	v1 := r.Group("/v1")
	{
		// registration
		v1.POST("/register", controller.RegisterNewUserViaEmailPassword)
	}
	// r.POST("/register/google", controller.registerNewUserUsingGoogle)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run(":8080")

	// ini untuk swagger
	// reference : https://golangexample.com/automatically-generate-restful-api-documentation-with-swagger-2-0-for-go/
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err != nil {
		log.Fatal(err)
	}
}

// Cors => allow access to non origin
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
