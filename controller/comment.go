package controller

import (
	"github.com/stevejo12/PMSFreelancer/models"
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserReviews => Add User Review
// GetUserReviews godoc
// @Summary Get User Review
// @Produce json
// @Accept  json
// @Tags User
// @Param token header string true "Token Header"
// @Success 200 {object} models.ResponseOKUserReviews
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /userReview [get]
func GetUserReviews(c *gin.Context) {
	id, err := strconv.Atoi(idToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something wrong with convertion string to int"})
		return
	}
	allData := []models.ReviewInfo{}

	result, err := config.DB.Query("SELECT * FROM comment WHERE user_id=? ORDER BY created_at ASC", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with query to get the reviews"})
		return
	}

	for result.Next() {
		data := models.ReviewDatabase{}
		returnData := models.ReviewInfo{}
		if err := result.Scan(&data.ID, &data.Message, &data.StarRating, &data.ProjectID, &data.MemberID, &data.UserID, &data.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		// get Project ID Detail Comment
		projectData, err := helpers.GetProjectDetailTitle(data.ProjectID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		// get the user who post the review
		reviewerData, err := helpers.GetUserInformationForReview(data.MemberID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		// userData, err := helpers.GetUserInformationForReview(id)

		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"code":    http.StatusInternalServerError,
		// 		"message": err.Error()})
		// 	return
		// }

		// organize the value
		returnData.ID = data.ID
		returnData.Message = data.Message
		returnData.StarRating = data.StarRating
		returnData.Project = projectData
		returnData.Reviewer = reviewerData
		// returnData.User = userData
		returnData.CreatedAt = data.CreatedAt

		allData = append(allData, returnData)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully Retrieved User Reviews",
		"data":    allData})
}

// AddUserReview => Add User Review
// AddUserReview godoc
// @Summary Add User Review
// @Produce json
// @Accept  json
// @Tags User
// @Param token header string true "Token Header"
// @Param Data body models.AddReview true "Data Format to add review"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /userReview [post]
func AddUserReview(c *gin.Context) {
	id := idToken
	param := models.AddReview{}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something wrong with convertion string to int"})
		return
	}

	err = c.BindJSON(&param)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	locationIndonesia, _ := time.LoadLocation("Asia/Jakarta")
	timeIndonesia := time.Now().In(locationIndonesia)

	formattedDate := fmt.Sprintf("%d-%02d-%02d", timeIndonesia.Year(), timeIndonesia.Month(), timeIndonesia.Day())

	query := "INSERT INTO comment(message, star_rating, project_id, member_id, user_id, created_at) VALUES"
	query = query + "(\"" + param.Message + "\", " + strconv.Itoa(param.StarRating) + ", " + strconv.Itoa(param.ProjectID) + ", " + id + ", " + strconv.Itoa(param.UserID) + ", \"" + formattedDate + "\")"

	_, err = config.DB.Exec(query)

	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully Added Review"})
}
