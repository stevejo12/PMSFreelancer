package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"
	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"
)

func getAllSkills() ([]models.UserSkills, error) {
	data, err := config.DB.Query("SELECT * FROM skills")
	defer data.Close()
	var allData []models.UserSkills

	if err != nil {
		return allData, errors.New("Server unable to execute query to database")
	}

	for data.Next() {
		// Scan one customer record
		var skills models.UserSkills
		if err := data.Scan(&skills.ID, &skills.Name, &skills.Created_at, &skills.Updated_at); err != nil {
			return []models.UserSkills{}, errors.New("Something is wrong with the database data")
		}
		allData = append(allData, skills)
	}
	if data.Err() != nil {
		return []models.UserSkills{}, errors.New("Something is wrong with the data retrieved")
	}

	return allData, nil
}

// GetAllSkills godoc
// @Summary Getting all list of skills
// @Produce json
// @Tags Skills
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /allSkills [get]
func GetAllSkills(c *gin.Context) {
	allSkills, err := getAllSkills()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    []models.UserSkills{}})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All Skills data have been retrieved",
		"data":    allSkills})
}

func getSkillNames(param string) ([]string, error) {
	var result []string

	if param == "" {
		return []string{}, nil
	}

	initialQuery, err := helpers.SettingInQueryWithID("skills", param, "*")

	if err != nil {
		return nil, err
	}

	data, err := config.DB.Query(initialQuery)
	defer data.Close()

	if err != nil {
		return result, errors.New("Server unable to execute query to database")
	}

	for data.Next() {
		// Scan one customer record
		var skills models.UserSkills
		if err := data.Scan(&skills.ID, &skills.Name, &skills.Created_at, &skills.Updated_at); err != nil {
			return []string{}, errors.New("Something is wrong with the database data")
		}
		result = append(result, skills.Name)
	}
	if data.Err() != nil {
		return []string{}, errors.New("Something is wrong with the data retrieved")
	}

	return result, nil
}

func userSkills(id string) ([]models.UserSkills, error) {
	query, err := helpers.SettingInQueryWithID("skills", id, "*")

	if err != nil {
		return []models.UserSkills{}, errors.New(err.Error())
	}

	resp, err := config.DB.Query(query)
	defer resp.Close()

	if err != nil {
		return []models.UserSkills{}, errors.New(err.Error())
	}

	allData := []models.UserSkills{}

	for resp.Next() {
		var databaseData models.UserSkills
		if err := resp.Scan(&databaseData.ID, &databaseData.Name, &databaseData.Created_at, &databaseData.Updated_at); err != nil {
			return []models.UserSkills{}, errors.New("Something is wrong with the database data")
		}

		allData = append(allData, databaseData)
	}

	if resp.Err() != nil {
		return []models.UserSkills{}, errors.New("Something is wrong with the data retrieved")
	}

	return allData, nil
}

// UpdateUserSkills => Updating User Skills
// UpdateUserSkills godoc
// @Summary Updating User Skills
// @Accept  json
// @Tags Skills
// @Param token header string true "Token Header"
// @Param Data body models.UpdateSkills true "New Skills ID List"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /editSkills [put]
func UpdateUserSkills(c *gin.Context) {
	id := idToken

	var data models.UpdateSkills

	err = c.BindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	strSkills := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(data.Skills)), ","), "[]")
	arrStrSkill := helpers.SplitComma(strSkills)

	var t2 = []int{}
	for _, i := range arrStrSkill {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		t2 = append(t2, j)
	}

	err = helpers.SkillList(t2)

	if err != nil {
		if err.Error() == "not exist" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "There is a value of skill that does not exist in the database id"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable execute query to the database"})
		return
	}

	_, err = config.DB.Exec("UPDATE login SET skill=? WHERE id=?", strSkills, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable execute query to the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All Skills data have been successfully updated"})
}
