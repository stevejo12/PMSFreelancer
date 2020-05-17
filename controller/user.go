package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"

	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var err error

var jwtKey = []byte("key_spirits")

type userPassword struct {
	Email    string
	Password string
}

type userGoogle struct {
	email    string
	googleID string
}

// RegisterUserWithPassword godoc
// @Summary Register new user using email and password
// @Produce json
// @Accept  json
// @Param account body userPassword true "Account"
// @Success 200 {object} models.ResponseWithNoBody
// @Router /register [post]
func RegisterUserWithPassword(c *gin.Context) {
	var newUser userPassword
	var multipleError []string

	err = c.Bind(&newUser)

	// checking empty body data
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

	// checking duplicate data
	var databaseUsername string

	err := config.DB.QueryRow("SELECT email FROM login WHERE email=?", newUser.Email).Scan(&databaseUsername)

	// this means email registered doesn't exist yet in the database.
	if err == sql.ErrNoRows {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusText(http.StatusInternalServerError),
				"message": "Server unable to hash the password into database"})

			return
		}

		// setting up data for inserting into database
		locationIndonesia, _ := time.LoadLocation("Asia/Jakarta")
		timeIndonesia := time.Now().In(locationIndonesia)

		formattedDate := fmt.Sprintf("%d-%02d-%02d", timeIndonesia.Year(), timeIndonesia.Month(), timeIndonesia.Day())

		selDB, err := config.DB.Prepare("INSERT INTO login(email, password, created_at, status) VALUES(?,?,?,?)")

		fmt.Println(err)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusText(http.StatusInternalServerError),
				"message": "Error preparing add user",
				"data":    []userPassword{}})

			return
		}

		// status at first created should be active
		_, err = selDB.Exec(newUser.Email, hashedPassword, formattedDate, "active")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusText(http.StatusInternalServerError),
				"message": "Server unable to execute query to database"})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusText(http.StatusOK),
			"message": "Adding user has been completed"})
	} else if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusText(http.StatusBadRequest),
			"message": "Email exists in the database"})
	} else {
		fmt.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Server unable to create your account"})
	}

}

// LoginUserWithPassword godoc
func LoginUserWithPassword(c *gin.Context) {
	var user userPassword

	err = c.Bind(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Error data format login"})

		return
	}

	var databaseEmail string
	var databasePassword string

	err = config.DB.QueryRow("SELECT email, password FROM login WHERE email=?", user.Email).Scan(&databaseEmail, &databasePassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Server unable to find the user email"})

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(user.Password))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Password doesn't match the email"})

		return
	}

	// setting token expiration
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.TokenClaims{
		Username: user.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusText(http.StatusOK),
		"message": "Login information is correct",
		"token":   jwtToken})
}

// ChangeUserPassword => Changing user password
func ChangeUserPassword(c *gin.Context) {
	var data models.ChangePassword
	var databaseEmail string
	var databasePassword string

	err := c.Bind(&data)

	// checking empty body data
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusText(http.StatusBadRequest),
			"message": "Data format is not as expected"})
		return
	}

	err = config.DB.QueryRow("SELECT email, password FROM login WHERE email=?", data.Email).Scan(&databaseEmail, &databasePassword)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Server unable to find the user email"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(data.OldPassword))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusText(http.StatusBadRequest),
			"message": "Old password is incorrect"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Server generate hashed password"})
		return
	}

	_, err = config.DB.Exec("UPDATE login SET password=? WHERE email=?", hashedPassword, data.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Server unable to execute query to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusText(http.StatusOK),
		"message": "Password has been successfully updated"})
}
