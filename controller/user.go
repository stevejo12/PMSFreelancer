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
	"gopkg.in/gomail.v2"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var err error

// CONFIG_SMTP_HOST => config hosting to send email
const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_EMAIL = "spirits.project.thesis@gmail.com"
const CONFIG_PASSWORD = "kevindjoni123"

var jwtKey = []byte("key_spirits")

type userPassword struct {
	Email    string
	Password string
}

// RegisterUserWithPassword godoc
// @Summary Register new user using email and password
// @Produce json
// @Accept  json
// @Param account body userPassword true "Account"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
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

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusText(http.StatusInternalServerError),
				"message": "Error preparing add user"})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Server unable to create your account"})
	}

}

// LoginUserWithPassword godoc
// @Summary Login user using email and password
// @Produce json
// @Accept  json
// @Param account body userPassword true "Account"
// @Success 200 {object} models.ResponseLoginWithToken
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /login [post]
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
// ChangeUserPassword godoc
// @Summary Register new user using email and password
// @Produce json
// @Accept  json
// @Param Info body models.ChangePassword true "Information needed to change password"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /change-password [put]
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

// HandleLogout => Log out from SPIRITS
func HandleLogout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusText(http.StatusOK),
		"message": "Successfully logged out"})
}

// ResetPassword => Sending email feature to reset password
func ResetPassword(c *gin.Context) {
	var data models.ResetPassword
	var mail string

	err := c.Bind(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusText(http.StatusBadRequest),
			"message": "Data format is not as expected"})
		return
	}

	err = config.DB.QueryRow("SELECT email FROM login WHERE email=?", data.Email).Scan(&mail)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusText(http.StatusBadRequest),
			"message": "Email is not registered in our database"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusBadRequest),
			"message": "Server unable to get information from database"})
		return
	}

	// generate link to the reset password
	// TO DO: make a generated link
	link := "https://localhost:8080/home"

	msg := fmt.Sprintf(`<html>
	<body>
	<p>Dear, asdf</p>
	<br>
	<p>You have requested a reset password for SPIRITS application</p>
	<p>This is the link to reset your password: %s</p>
	<p>Note: This link expires in 30 minutes after this email is recieved</p>
	<br>
	<p>Best Regards,</p>
	SPIRITS Team
	</body>
	</html>`, link)

	// setting up the content for the email
	m := gomail.NewMessage()
	m.SetHeader("From", CONFIG_EMAIL)
	m.SetHeader("To", mail)
	m.SetHeader("Subject", "You have requested to reset password")
	m.SetBody("text/html", msg)

	// Send the email to user
	d := gomail.NewPlainDialer(CONFIG_SMTP_HOST, CONFIG_SMTP_PORT, CONFIG_EMAIL, CONFIG_PASSWORD)
	if err := d.DialAndSend(m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusText(http.StatusInternalServerError),
			"message": "Server unable to send email to user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusText(http.StatusOK),
		"message": "Email has been send"})
}
