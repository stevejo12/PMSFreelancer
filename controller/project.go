package controller

import (
	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"

	"github.com/gin-gonic/gin"
)

// SearchProject => Search project in SPIRITS
func SearchProject(c *gin.Context) {
	// initialize variables
	// page is page number in pagination
	// size is the number of result per page
	pageParam, ok := c.Request.URL.Query()["page"]
	sizeParam, ok := c.Request.URL.Query()["size"]

	if !ok || len(pageParam[0]) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Request Url should have page and size in it"})
		return
	}

	page, err := strconv.Atoi(pageParam[0])

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Page value is not integer"})
		return
	}

	size, err := strconv.Atoi(sizeParam[0])

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Size value is not integer"})
		return
	}

	// record #1 is number 0
	var startingRecordNumber = page * size
	var endingRecordNumber = startingRecordNumber + size

	result, err := config.DB.Query("SELECT id, title, description, price FROM project ORDER BY ID ASC LIMIT ?,?", startingRecordNumber, endingRecordNumber)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with query to get the project list"})
		return
	}

	var allData []models.SearchProjectQuery
	for result.Next() {
		var project models.SearchProjectQuery
		if err := result.Scan(&project.ID, &project.Title, &project.Description, &project.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}
		allData = append(allData, project)
	}

	if result.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with the data retrieved"})
		return
	}

	var resp models.SearchProjectResponse

	resp.Project = allData
	resp.TotalSearch = len(allData)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All Project data have been retrieved",
		"data":    allData})
}

func getProjectLinks(param string) ([]string, error) {
	var result []string
	query, err := helpers.SettingInQueryWithID("project_links", param)

	if err != nil {
		return nil, err
	}

	data, err := config.DB.Query(query)

	if err != nil {
		return nil, errors.New("Server unable to execute query to database")
	}

	for data.Next() {
		// Scan one customer record
		var link models.ProjectLinksResponse
		if err := data.Scan(&link.ID, &link.Project_link); err != nil {
			return []string{}, errors.New("Something is wrong with the database data")
		}
		result = append(result, link.Project_link)
	}
	if data.Err() != nil {
		return []string{}, errors.New("Something is wrong with the data retrieved")
	}

	return result, nil
}

func AddProject(c *gin.Context) {
	id := c.Param("id")

	var param models.CreateProject

	err := c.Bind(&param)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	_, err = config.DB.Query("INSERT INTO project(title, description, skills, price, attachment, owner_id) VALUES(?,?,?,?,?,?)", param.Title, param.Description, param.Skills, param.Price, param.Attachment, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully Added New Project"})
}

func GetAllUserProjects(c *gin.Context) {
	// user id
	id := c.Param("id")

	result, err := config.DB.Query("SELECT id, title, description, status FROM project WHERE owner_id=?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with query to get the project list"})
		return
	}

	var allData []models.GetUserProjectResponse

	for result.Next() {
		var project models.GetUserProjectResponse

		if err = result.Scan(&project.ID, &project.Title, &project.Description, &project.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		allData = append(allData, project)
	}

	if result.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with the data retrieved"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All Project data have been retrieved",
		"data":    allData})
}

func ProjectDetail(c *gin.Context) {
	// project id
	id := c.Param("id")

	result, err := config.DB.Query("SELECT id, title, skills, price, attachment, owner_id, interested_members FROM project WHERE id=?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with query to get the project detail"})
		return
	}

	var allData []models.ProjectDetailResponse
	for result.Next() {
		var dbResult models.ProjectDetailRequest
		var data models.ProjectDetailResponse

		if err = result.Scan(&dbResult.ID, &dbResult.Title, &dbResult.Skills, &dbResult.Price, &dbResult.Attachment, &dbResult.OwnerID, &dbResult.InterestedMembers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		// get the detail about the skills
		dataSkills, err := getSkillNames(dbResult.Skills)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err})
			return
		}

		// get the detail link of the attachment for this project
		var dataLink []string
		if dbResult.Attachment.Valid {
			dataLink, err = getProjectLinks(dbResult.Attachment.String)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Something is wrong with the database data"})
				return
			}
		} else {
			dataLink = []string{}
		}

		// get the detail about the owner
		var ownerInfo models.OwnerInfo
		ownerData, err := config.DB.Query("SELECT id, first_name, last_name, location, created_at FROM login WHERE id=?", dbResult.OwnerID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		for ownerData.Next() {
			var queryResult models.OwnerInfoQuery
			if err := ownerData.Scan(&queryResult.ID, &queryResult.FirstName, &queryResult.LastName, &queryResult.Location, &queryResult.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Something is wrong with the database data"})
				return
			}

			ownerInfo.ID = queryResult.ID
			ownerInfo.FullName = queryResult.FirstName + " " + queryResult.LastName
			ownerInfo.Location = queryResult.Location
			memberInfo := helpers.SplitDash(queryResult.CreatedAt)
			if len(memberInfo) == 3 {
				ownerInfo.Member = memberInfo[0]
			} else {
				ownerInfo.Member = ""
			}
		}

		// get # of completed project
		var count int
		err = config.DB.QueryRow("SELECT COUNT(*) FROM project WHERE owner_id=?", id).Scan(&count)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server unable to execute query to database"})
			return
		}
		ownerInfo.ProjectCompleted = count

		// construct the response to user
		data.ID = dbResult.ID
		data.Title = dbResult.Title
		data.Skills = dataSkills
		data.Price = dbResult.Price
		data.Owner = ownerInfo
		data.Attachment = dataLink

		allData = append(allData, data)
	}

	if len(allData) != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Project data have been retrieved",
		"data":    allData})
}

// SubmitProjectInterest => Potential Freelancer submit their interest before accepted by project owner
// /:id => to get the project ID
// parameter ID: this is the freelancer id to register
func SubmitProjectInterest(c *gin.Context) {
	// this is project id
	id := c.Param("id")

	// this ID is for the potential freelancer id
	type submitInterest struct {
		ID int
	}

	var param submitInterest

	err := c.Bind(&param)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// check if this is the owner
	owner, err := helpers.IsThisIDProjectOwner(id, param.ID)

	if !owner {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		// check if this member already registered
		member, err := helpers.IsThisMemberRegistered(id, param.ID)

		if !member {
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": err.Error()})
				return
			}

			// if they are not member yet, register them
			ok := registerUserToInterestedMembers(id, param.ID)

			if ok != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": ok.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "Member has been added to interested list"})
		} else {
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": err.Error()})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something wrong with the server"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "This user is the project owner"})
		return
	}
}

func registerUserToInterestedMembers(projectID string, userID int) error {
	listMember, err := helpers.GetMemberList(projectID)

	if err != nil {
		return err
	}

	// add member to the list
	if listMember == "" {
		listMember = strconv.Itoa(userID)
	} else {
		listMember = listMember + "," + strconv.Itoa(userID)
	}

	_, err = config.DB.Exec("UPDATE project SET interested_members=? WHERE id=?", listMember, projectID)

	if err != nil {
		return errors.New("Server is unable to execute query to the database")
	}

	return nil
}
