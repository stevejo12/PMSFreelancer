package controller

import (
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"

	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func getUserReviews(id string) ([]models.ReviewInfo, error) {
	allData := []models.ReviewInfo{}

	result, err := config.DB.Query("SELECT * FROM comment WHERE user_id=? ORDER BY created_at DESC", id)

	if err != nil {
		return []models.ReviewInfo{}, errors.New("Something is wrong with query to get the reviews")
	}

	for result.Next() {
		returnData := models.ReviewInfo{}
		var isOwner, projectID, userID, memberID int
		var dateCreated string
		if err := result.Scan(&returnData.ID, &returnData.Message, &returnData.StarRating, &projectID, &memberID, &userID, &dateCreated, &isOwner); err != nil {
			return []models.ReviewInfo{}, errors.New("Something is wrong with the database data")
		}

		// get Project ID Detail Comment
		projectData, err := helpers.GetProjectDetailTitle(projectID)

		if err != nil {
			return []models.ReviewInfo{}, err
		}

		// get the user who post the review
		reviewerData, err := helpers.GetUserInformationForReview(memberID)

		if err != nil {
			return []models.ReviewInfo{}, err
		}

		if isOwner == 0 {
			reviewerData.IsOwner = false
		} else {
			reviewerData.IsOwner = true
		}

		if dateCreated != "" {
			arrDate := strings.Split(dateCreated, "T")

			if len(arrDate) > 0 {
				returnData.CreatedAt = arrDate[0]
			} else {
				returnData.CreatedAt = ""
			}
		} else {
			returnData.CreatedAt = ""
		}

		// organize the value
		returnData.Project = projectData
		returnData.Reviewer = reviewerData

		allData = append(allData, returnData)
	}

	return allData, nil
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
	id, err := strconv.Atoi(idToken)
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

	var ownerID, freelancerID int
	err = config.DB.QueryRow("SELECT owner_id, accepted_memberid FROM project WHERE id=?", param.ProjectID).Scan(&ownerID, &freelancerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to get owner and freelancer information"})
		return
	}

	if ownerID != id && freelancerID != id {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "This user is neither owner nor freelancer on the project"})
		return
	}

	locationIndonesia, _ := time.LoadLocation("Asia/Jakarta")
	timeIndonesia := time.Now().In(locationIndonesia)

	formattedDate := fmt.Sprintf("%d-%02d-%02d", timeIndonesia.Year(), timeIndonesia.Month(), timeIndonesia.Day())

	var writer, recipient, isOwnerVal int
	if param.IsOwner {
		writer = ownerID
		recipient = freelancerID
		isOwnerVal = 1
	} else {
		writer = freelancerID
		recipient = ownerID
		isOwnerVal = 0
	}

	// member_id => the one who writes the comment
	// user_id => the reviewee
	query := "INSERT INTO comment(message, star_rating, project_id, member_id, user_id, created_at, is_owner) VALUES"
	query = query + "(\"" + param.Message + "\", " + strconv.Itoa(param.StarRating) + ", " + strconv.Itoa(param.ProjectID) + ", " + strconv.Itoa(writer) + ", " + strconv.Itoa(recipient) + ", \"" + formattedDate + "\", " + strconv.Itoa(isOwnerVal) + ")"

	_, err = config.DB.Exec(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully Added Review"})
}
