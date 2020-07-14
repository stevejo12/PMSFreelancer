package helpers

import (
	// "PMSFreelancer/models"
	"math"
	"github.com/stevejo12/PMSFreelancer/models"
)

// GetAverageUserRating => get user average rating
func GetAverageUserRating(reviews []models.ReviewInfo) (float64, error) {
	avgRating := 0.0
	if len(reviews) == 0 {
		return avgRating, nil
	}

	var totalRating int
	totalData := len(reviews)

	for i := 0; i < len(reviews); i++ {
		rating := reviews[i].StarRating

		totalRating += rating
	}

	// get the average
	avgRating = float64(totalRating) / float64(totalData)

	// change to 1 decimal place
	avgRating = math.Round(avgRating*10) / 10

	return avgRating, nil

}
