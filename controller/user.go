package controller

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"

	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"

	"github.com/gin-gonic/gin"
)

var err error

// CONFIG_SMTP_HOST => config hosting to send email
const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_EMAIL = "spirits.project.thesis@gmail.com"
const CONFIG_PASSWORD = "kevindjoni123"

var jwtKey = []byte("key_spirits")

// RegisterUserWithPassword godoc
// @Summary Register new user using email and password
// @Produce json
// @Accept  json
// @Tags User
// @Param account body models.RegistrationUserUsingPassword true "Account"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /register [post]
func RegisterUserWithPassword(c *gin.Context) {
	var newUser models.RegistrationUserUsingPassword
	var multipleError []string

	err = c.BindJSON(&newUser)

	// checking empty body data
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Error binding new user"})

		return
	}

	ref := reflect.ValueOf(&newUser).Elem()

	for i := 0; i < ref.NumField(); i++ {
		varName := ref.Type().Field(i).Name
		// varType := ref.Type().Field(i).Type
		varValue := ref.Field(i).Interface()

		strVal := helpers.ConvertToString(varValue)

		if strVal == "" || (strVal == "0" && varName == "Location") {
			message := varName + " must not be empty"
			multipleError = append(multipleError, message)
		}
	}

	if len(multipleError) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": multipleError})
		return
	}

	// checking duplicate data
	emailExist := helpers.CheckDuplicateEmail(newUser.Email)
	usernameExist := helpers.CheckDuplicateUsername(newUser.Username)

	// this means email registered and username registered don't exist yet in the database.
	if emailExist == nil && usernameExist == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server unable to hash the password into database"})

			return
		}

		// checking full name
		// splitting into first name and last name
		firstname, lastname, err := helpers.SplittingFullname(newUser.Fullname)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Fullname is empty"})
			return
		}

		// checking the skills list provided
		err = helpers.SkillList(newUser.Skills)

		if err != nil {
			if err.Error() == "not exist" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "There is a value that does not exist in the database id"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		// checking location (expect the country id)
		err = helpers.CountryList(newUser.Location)

		if err != nil {
			if err.Error() == "not exist" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "Country ID does not exist in the database"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		// setting up data for inserting into database
		locationIndonesia, _ := time.LoadLocation("Asia/Jakarta")
		timeIndonesia := time.Now().In(locationIndonesia)

		formattedDate := fmt.Sprintf("%d-%02d-%02d", timeIndonesia.Year(), timeIndonesia.Month(), timeIndonesia.Day())

		selDB, err := config.DB.Prepare("INSERT INTO login(email, password, created_at, status, username, description,first_name, last_name, location, skill) VALUES(?,?,?,?,?,?,?,?,?,?)")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Error preparing add user"})
			return
		}

		// status at first created should be active
		_, err = selDB.Exec(newUser.Email, hashedPassword, formattedDate, "active", newUser.Username, newUser.Description, firstname, lastname, newUser.Location, strings.Join(newUser.Skills, ","))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server unable to execute query to database"})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Adding user has been completed"})
	} else if emailExist != nil || usernameExist != nil {
		multipleError = []string{}
		if emailExist != nil {
			multipleError = append(multipleError, emailExist.Error())
		}
		if usernameExist != nil {
			multipleError = append(multipleError, usernameExist.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": multipleError})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to create your account"})
		return
	}
}

func RegisterUserWithGoogle(c *gin.Context) {
	var newUser models.RegistrationUserUsingGoogle
	var multipleError []string

	err = c.BindJSON(&newUser)

	// checking empty body data
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Error binding new user"})

		return
	}

	ref := reflect.ValueOf(&newUser).Elem()

	for i := 0; i < ref.NumField(); i++ {
		varName := ref.Type().Field(i).Name
		// varType := ref.Type().Field(i).Type
		varValue := ref.Field(i).Interface()

		strVal := helpers.ConvertToString(varValue)

		if strVal == "" || (strVal == "0" && varName == "Location") {
			message := varName + " must not be empty"
			multipleError = append(multipleError, message)
		}
	}

	if len(multipleError) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": multipleError})
		return
	}

	// checking duplicate data
	emailExist := helpers.CheckDuplicateEmail(newUser.Email)
	usernameExist := helpers.CheckDuplicateUsername(newUser.Username)

	// this means email registered and username registered don't exist yet in the database.
	if emailExist == nil && usernameExist == nil {
		// checking full name
		// splitting into first name and last name
		firstname, lastname, err := helpers.SplittingFullname(newUser.Fullname)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Fullname is empty"})
			return
		}

		// checking the skills list provided
		err = helpers.SkillList(newUser.Skills)

		if err != nil {
			if err.Error() == "not exist" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "There is a value that does not exist in the database id"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		// checking location (expect the country id)
		err = helpers.CountryList(newUser.Location)

		if err != nil {
			if err.Error() == "not exist" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "Country ID does not exist in the database"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		// setting up data for inserting into database
		locationIndonesia, _ := time.LoadLocation("Asia/Jakarta")
		timeIndonesia := time.Now().In(locationIndonesia)

		formattedDate := fmt.Sprintf("%d-%02d-%02d", timeIndonesia.Year(), timeIndonesia.Month(), timeIndonesia.Day())

		selDB, err := config.DB.Prepare("INSERT INTO login(email, google_id, created_at, status, username, description,first_name, last_name, location, skill) VALUES(?,?,?,?,?,?,?,?,?,?)")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Error preparing add user"})
			return
		}

		// status at first created should be active
		_, err = selDB.Exec(newUser.Email, newUser.GoogleID, formattedDate, "active", newUser.Username, newUser.Description, firstname, lastname, newUser.Location, strings.Join(newUser.Skills, ","))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server unable to execute query to database"})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Adding user has been completed"})
	} else if emailExist != nil || usernameExist != nil {
		multipleError = []string{}
		if emailExist != nil {
			multipleError = append(multipleError, emailExist.Error())
		}
		if usernameExist != nil {
			multipleError = append(multipleError, usernameExist.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": multipleError})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to create your account"})
		return
	}
}

// LoginUserWithPassword godoc
// @Summary Login user using email and password
// @Produce json
// @Accept  json
// @Tags User
// @Param parameter body models.LoginUserPassword true "Account"
// @Success 200 {object} models.ResponseOKLogin
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /login [post]
func LoginUserWithPassword(c *gin.Context) {
	var user models.LoginUserPassword

	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		fmt.Println("err", err.Error())
	}

	fmt.Println("body of the request: ", string(body))

	c.Request.Body = ioutil.NopCloser(bytes.NewReader([]byte(body)))

	err = c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Error data format login"})
		return
	}

	// tambahin ERROR DEH KALO USER EMAIL / PASSWORD KOSONG

	fmt.Println(user)
	fmt.Println("user email", user.Email)
	fmt.Println("user pass", user.Password)

	var databaseID string
	var databaseEmail string
	var databasePassword string

	err = config.DB.QueryRow("SELECT id, email, password FROM login WHERE email=?", &user.Email).Scan(&databaseID, &databaseEmail, &databasePassword)

	if err != nil {
		fmt.Println("err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to find the user email"})

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(user.Password))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Password doesn't match the email"})

		return
	}

	// setting token expiration
	cookieToken, expirationTime, err := generateToken(databaseID)

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

// ChangeUserPassword => Changing user password
// ChangeUserPassword godoc
// @Summary Register new user using email and password
// @Produce json
// @Accept  json
// @Tags User
// @Param Info body models.ChangePassword true "Information needed to change password"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /change-password [put]
func ChangeUserPassword(c *gin.Context) {
	var data models.ChangePassword
	var databaseEmail string
	var databasePassword string

	err := c.BindJSON(&data)

	// checking empty body data
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is not as expected"})
		return
	}

	err = config.DB.QueryRow("SELECT email, password FROM login WHERE email=?", data.Email).Scan(&databaseEmail, &databasePassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to find the user email"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(data.OldPassword))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Old password is incorrect"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server generate hashed password"})
		return
	}

	_, err = config.DB.Exec("UPDATE login SET password=? WHERE email=?", hashedPassword, data.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Password has been successfully updated"})
}

// HandleLogout => Log out from SPIRITS
// HandleLogout godoc
// @Summary Logout
// @Produce json
// @Tags User
// @Param token header string true "Token Header"
// @Success 200 {object} models.ResponseWithNoBody
// @Router /logout [post]
func HandleLogout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully logged out"})
}

// ResetPassword => Sending email feature to reset password
// ResetPassword godoc
// @Summary Reset Password (Forget Password Feature)
// @Produce json
// @Accept  json
// @Tags User
// @Param Info body models.ResetPassword true "Information needed to change password"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /resetPassword [post]
func ResetPassword(c *gin.Context) {
	var data models.ResetPassword
	var mail string

	err := c.BindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is not as expected"})
		return
	}

	err = config.DB.QueryRow("SELECT email FROM login WHERE email=?", data.Email).Scan(&mail)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Email is not registered in our database"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Server unable to get information from database"})
		return
	}

	// generate link to the reset password
	// TO DO: update this link from localhost to real url for resetting password
	link := "https://localhost:8080/reset-password/"

	// generate token for link
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(mail), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to do hashing"})
		return
	}

	// generate template to insert into database
	// 30 mins lifespan time for the link
	lifetime := 30 * 60 * time.Second
	locationIndonesia, _ := time.LoadLocation("Asia/Jakarta")
	timeNow := time.Now().In(locationIndonesia)

	timeNow = timeNow.Add(lifetime)

	_, err = config.DB.Exec("INSERT INTO resetpassword_token(email, token, expire) VALUES(?,?,?)", mail, string(hashedToken), timeNow)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to store token into databse"})
		return
	}

	hashStringVal := base64.StdEncoding.EncodeToString(hashedToken)

	link = link + "reset_token=" + hashStringVal

	msg := fmt.Sprintf(`<html>
	<body>
	<p>We received request to reset your password</p>
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
			"code":    http.StatusInternalServerError,
			"message": "Server unable to send email to user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Email has been send"})
}

// UpdateNewPassword => Updating user password without old password
// ResetPassword godoc
// @Summary Update Password (after reset password)
// @Produce json
// @Accept  json
// @Tags User
// @Param Info body models.UpdateResetPassword true "Information needed to update password"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /updateNewPassword [post]
func UpdateNewPassword(c *gin.Context) {
	var param models.UpdateResetPassword

	err = c.BindJSON(&param)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	if param.Token != "" {
		var dbData models.DatabaseResetPassword

		tokenValue, err := base64.StdEncoding.DecodeString(param.Token)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server unable to decode the token"})
			return
		}

		stringTokenValue := string(tokenValue)

		err = config.DB.QueryRow("SELECT email, token, expire FROM resetpassword_token WHERE token=?", stringTokenValue).Scan(&dbData.Email, &dbData.Token, &dbData.Expire)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server unable to find the token in the database"})
			return
		}

		timeNow := time.Now()
		parseTime, err := time.Parse("2006-06-05 00:00:00", dbData.Expire)

		// check if
		expire := timeNow.Before(parseTime)

		if expire {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Token for this link has expired"})
			return
		}

		// hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server unable to hash the password into database"})
			return
		}

		_, err = config.DB.Exec("UPDATE login SET password=? WHERE email=?", hashedPassword, dbData.Email)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server unable to update the password into database"})
			return
		}

		// remove token to avoid spam update after successfully doing it once (in the database)
		_, err = config.DB.Exec("DELETE FROM resetpassword_token WHERE token=?", stringTokenValue)

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "User password has been successfully updated"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Parameter Token is empty"})
		return
	}
}

// GetUserProfile => Updating user password without old password
// GetUserProfile godoc
// @Summary User Profile Data
// @Produce json
// @Accept  json
// @Tags User
// @Param token header string true "Token Header"
// @Success 200 {object} models.ResponseOKGetUserProfile
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /userProfile [get]
func GetUserProfile(c *gin.Context) {
	id := idToken
	var data models.UserProfile

	// get education list
	educationData, err := userEducation(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// get experience list
	experienceData, err := userExperience(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	var dataQuery models.QueryUserProfile
	var picData sql.NullString

	err = config.DB.QueryRow("SELECT id, first_name, last_name, email, description, picture, created_at, username, location, skill, balance FROM login WHERE id=?", id).Scan(&dataQuery.ID, &dataQuery.Firstname, &dataQuery.LastName, &dataQuery.Email, &dataQuery.Description, &picData, &dataQuery.CreatedAt, &dataQuery.Username, &dataQuery.Location, &dataQuery.Skills, &dataQuery.Balance)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to the database"})
		return
	}

	// get skill list
	skillData, err := userSkills(dataQuery.Skills)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// get user portfolio
	userPortfolio, err := allUserPortfolio(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// check if there is any uploaded picture
	// valid means there is value
	// else means null
	if picData.Valid {
		dataQuery.Picture = picData.String
	} else {
		dataQuery.Picture = ""
	}

	// Year of Member
	arrMemberSince := helpers.SplitDash(dataQuery.CreatedAt)
	if len(arrMemberSince) > 0 {
		data.Member = arrMemberSince[0]
	} else {
		data.Member = ""
	}

	// arrange all data
	data.Education = educationData
	data.Experience = experienceData
	data.Skill = skillData
	data.ID = dataQuery.ID
	// data.Fullname = dataQuery.Firstname + " " + dataQuery.LastName
	data.FirstName = dataQuery.Firstname
	data.LastName = dataQuery.LastName
	data.Email = dataQuery.Email
	data.Description = dataQuery.Description
	data.Picture = dataQuery.Picture
	data.Username = dataQuery.Username
	data.Location = dataQuery.Location
	data.Portfolio = userPortfolio
	data.Balance = dataQuery.Balance

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All User Profile data have been successfully retrieved",
		"data":    data})
}
