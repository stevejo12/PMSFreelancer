package controller

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/stevejo12/PMSFreelancer/helpers"

	"github.com/gin-gonic/gin"
)

var err error

type newUserPasswordRegistration struct {
	Email    string
	Password string
}

type newUserGoogleRegistration struct {
	email    string
	googleID string
}

// RegisterNewUserViaEmailPassword godoc
// @Summary Register new user using email and password
// @Produce json
// @Accept  json
// @Param account body newUserPasswordRegistration true "Account"
// @Success 200 {object} models.ResponseWithNoBody
// @Router /register [post]
func RegisterNewUserViaEmailPassword(c *gin.Context) {
	var newUser newUserPasswordRegistration
	var multipleError []string

	err = c.Bind(&newUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Error binding new user"})

		return
	}

	ref := reflect.ValueOf(&newUser).Elem()

	for i := 0; i < ref.NumField(); i++ {
		varName := ref.Type().Field(i).Name
		// varType := ref.Type().Field(i).Type
		varValue := ref.Field(i).Interface()

		strVal := helpers.ConvertToString(varValue)

		if strVal == "" {
			message := varName + " must not be empty"
			multipleError = append(multipleError, message)
		}
	}

	if len(multipleError) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": multipleError})

		return
	}

	locationIndonesia, _ := time.LoadLocation("Asia/Jakarta")
	timeIndonesia := time.Now().In(locationIndonesia)

	formattedDate := fmt.Sprintf("%d-%02d-%02d", timeIndonesia.Year(), timeIndonesia.Month(), timeIndonesia.Day())

	selDB, err := DB.Prepare("INSERT INTO login(email, password, created_at, status) VALUES(?,?,?,?)")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Error preparing add user",
			"data":    []newUserPasswordRegistration{}})

		return
	}

	// status at first created should be active
	_, err = selDB.Exec(newUser.Email, newUser.Password, formattedDate, "active")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": err})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusText(http.StatusOK),
		"message": "Adding user has been completed"})
}
