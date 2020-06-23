package controller

import (
	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/helpers"
	"github.com/stevejo12/PMSFreelancer/models"

	"github.com/gin-gonic/gin"
)

func getProjectAttachments(projectID string) ([]models.ProjectLinksResponse, error) {
	result := []models.ProjectLinksResponse{}

	data, err := config.DB.Query("SELECT id, project_link FROM project_links WHERE project_id=?", projectID)

	if err != nil {
		return nil, errors.New("Server unable to execute query to database")
	}

	for data.Next() {
		// Scan one customer record
		var link models.ProjectLinksResponse
		if err := data.Scan(&link.ID, &link.Project_link); err != nil {
			return []models.ProjectLinksResponse{}, errors.New("Something is wrong with the database data")
		}
		result = append(result, link)
	}
	if data.Err() != nil {
		return []models.ProjectLinksResponse{}, errors.New("Something is wrong with the data retrieved")
	}

	return result, nil
}

// AddProject => Add User Education
// AddProject godoc
// @Summary Adding User Project
// @Produce json
// @Accept  json
// @Tags Project
// @Param token header string true "Token Header"
// @Param Data body models.CreateProject true "Data Format to add education"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /addProject [post]
func AddProject(c *gin.Context) {
	id := idToken

	var param models.CreateProject

	err := c.BindJSON(&param)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	err = helpers.SkillList(param.Skills)

	if err != nil {
		if err.Error() == "not exist" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "There is skill id does not exist in the database id"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	skillDataQuery := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(param.Skills)), ","), "[]")

	queryResult, err := config.DB.Exec("INSERT INTO project(title, description, skills, price, owner_id) VALUES(?,?,?,?,?)", param.Title, param.Description, skillDataQuery, param.Price, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})
		return
	}

	rowID, err := queryResult.LastInsertId()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to retrieve the id inserted"})
		return
	}

	// add project link list
	for i := 0; i < len(param.Attachment); i++ {
		_, err = config.DB.Exec("INSERT INTO project_links(project_link, project_id) VALUES(?,?)", param.Attachment[i], rowID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server unable to execute query to database"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully Added New Project"})
}

// GetAllUserProjects => List of User Projects
// GetAllUserProjects godoc
// @Summary User Projects
// @Produce json
// @Accept  json
// @Tags Project
// @Param token header string true "Token Header"
// @Success 200 {object} models.ResponseOKGetUserProject
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /userProjects [get]
func GetAllUserProjects(c *gin.Context) {
	// user id
	id := idToken

	result, err := config.DB.Query("SELECT id, title, description, status FROM project WHERE owner_id=?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with query to get the project list"})
		return
	}

	allData := []models.GetUserProjectResponse{}

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

// ProjectDetail => Project Detail
// ProjectDetail godoc
// @Summary User Project Detail
// @Produce json
// @Accept  json
// @Tags Project
// @Param id path int64 true "Project ID"
// @Success 200 {object} models.ResponseOKProjectDetail
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /projectDetail/{id} [get]
func ProjectDetail(c *gin.Context) {
	// project id
	id := c.Param("id")

	result, err := config.DB.Query("SELECT id, title, skills, price, owner_id, interested_members FROM project WHERE id=?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with query to get the project detail"})
		return
	}

	allData := []models.ProjectDetailResponse{}
	for result.Next() {
		var dbResult models.ProjectDetailRequest
		var data models.ProjectDetailResponse

		if err = result.Scan(&dbResult.ID, &dbResult.Title, &dbResult.Skills, &dbResult.Price, &dbResult.OwnerID, &dbResult.InterestedMembers); err != nil {
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
				"message": err.Error()})
			return
		}

		// get the detail link of the attachment for this project
		dataLink, err := getProjectAttachments(id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
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
			ownerInfo.FirstName = queryResult.FirstName
			ownerInfo.LastName = queryResult.LastName
			ownerInfo.Location = queryResult.Location
			memberInfo := helpers.SplitDash(queryResult.CreatedAt)
			if len(memberInfo) == 3 {
				ownerInfo.Member = memberInfo[0]
			} else {
				ownerInfo.Member = ""
			}
		}

		// get # of completed project
		projectCompleted, err := helpers.GetUserCompletedProject(id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}
		ownerInfo.ProjectCompleted = projectCompleted

		// get the names and id of the member
		interestedMembers, err := helpers.GetInterestedMemberNames(id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		// construct the response to user
		data.ID = dbResult.ID
		data.Title = dbResult.Title
		data.Skills = dataSkills
		data.Price = dbResult.Price
		data.Owner = ownerInfo
		data.Attachment = dataLink
		data.InterestedMembers = interestedMembers

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
// SubmitProjectInterest godoc
// @Summary Submit Project Interest
// @Accept  json
// @Tags Project
// @Param token header string true "Token Header"
// @Param id path int64 true "Project ID"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /submitProjectInterest/{id} [post]
func SubmitProjectInterest(c *gin.Context) {
	// this is project id
	id := c.Param("id")

	userID, err := strconv.Atoi(idToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Unable to retrieve user id"})
		return
	}

	// check if this is the owner
	owner, err := helpers.IsThisIDProjectOwner(id, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	if !owner {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		// check if this member already registered
		member, err := helpers.IsThisMemberRegistered(id, userID)

		if !member {
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": err.Error()})
				return
			}

			// if they are not member yet, register them
			ok := registerUserToInterestedMembers(id, userID)

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
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "This member has already been registered"})
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

// AcceptProjectInterest godoc
// @Summary Accepting Freelancer to Project
// @Accept  json
// @Tags Project
// @Param token header string true "Token Header"
// @Param id path int64 true "Project ID"
// @Param data body models.ProjectAcceptMemberParameter true "Freelancer which you want to accept"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /acceptProjectInterest/{id} [post]
func AcceptProjectInterest(c *gin.Context) {
	id := c.Param("id")

	ownerID, err := strconv.Atoi(idToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something wrong with convertion string to int"})
		return
	}

	param := models.ProjectAcceptMemberParameter{}

	err = c.BindJSON(&param)

	if err != nil || param.FreelancerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	// check if status is currently Listed
	var status string
	err = config.DB.QueryRow("SELECT status FROM project WHERE id=?", id).Scan(&status)

	if status != "Listed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "This function works only when the project status is Listed"})
		return
	}

	// check if the owner id is correct
	owner, err := helpers.IsThisIDProjectOwner(id, ownerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	if owner {
		// check if the member is registered
		member, err := helpers.IsThisMemberRegistered(id, param.FreelancerID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		if member {
			// Create Trello Board and store the link
			var boardTitle string
			err = config.DB.QueryRow("SELECT title FROM project WHERE id=?", id).Scan(&boardTitle)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Server is unable to execute query to the database"})
				return
			}

			trelloBoardID, err := createTrelloBoard(boardTitle, param.TrelloKey)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": err.Error()})
				return
			}

			trelloURL := "https://www.trello.com/b/" + trelloBoardID

			// Update the project data
			_, err = config.DB.Exec("UPDATE project SET status=?, accepted_memberid=?, trello_url=? WHERE id=?", "On Going", param.FreelancerID, trelloURL, id)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Server is unable to execute query to the database"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "Successfully accepted member"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "This member ID is not listed yet as interested member"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "This Owner doesn't own the project"})
		return
	}
}

// ReviewProject godoc
// @Summary Submit Project for review
// @Accept  json
// @Tags Project
// @Param token header string true "Token Header"
// @Param id path int64 true "Project ID"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /submitProjectForReview/{id} [post]
func ReviewProject(c *gin.Context) {
	// project id
	id := c.Param("id")

	freelancerID, err := strconv.Atoi(idToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// check if status is currently Listed
	var status string
	err = config.DB.QueryRow("SELECT status FROM project WHERE id=?", id).Scan(&status)

	if status != "On Going" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "The status isn't on going yet"})
		return
	}

	freelancer, err := helpers.IsThisTheAcceptedMember(id, freelancerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	if freelancer {
		_, err = config.DB.Exec("UPDATE project SET status=? WHERE id=?", "On Review", id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server is unable to execute query to the database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Successfully updating project to on review"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This user is unauthorized to make the project to be reviewed"})
		return
	}
}

// CompleteProject godoc
// @Summary Review Done For Project Owner
// @Accept  json
// @Tags Project
// @Param token header string true "Token Header"
// @Param id path int64 true "Project ID"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /completeProject/{id} [post]
func CompleteProject(c *gin.Context) {
	// project id
	id := c.Param("id")

	ownerID, err := strconv.Atoi(idToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// check if status is currently On Review
	var status string
	err = config.DB.QueryRow("SELECT status FROM project WHERE id=?", id).Scan(&status)

	if status != "On Review" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "The status isn't on review yet"})
		return
	}

	owner, err := helpers.IsThisIDProjectOwner(id, ownerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	if owner {
		_, err = config.DB.Exec("UPDATE project SET status=? WHERE id=?", "Done", id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server is unable to execute query to the database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Successfully updating project to Done"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This user is unauthorized to make the project DONE"})
		return
	}
}

func RejectReviewProject(c *gin.Context) {
	// project id
	id := c.Param("id")

	ownerID, err := strconv.Atoi(idToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// check if status is currently On Review
	var status string
	err = config.DB.QueryRow("SELECT status FROM project WHERE id=?", id).Scan(&status)

	if status != "On Review" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "The status isn't on review yet"})
		return
	}

	owner, err := helpers.IsThisIDProjectOwner(id, ownerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	if owner {
		_, err = config.DB.Exec("UPDATE project SET status=? WHERE id=?", "On Going", id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server is unable to execute query to the database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Successfully updating project to On Going"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "This user is unauthorized to reject the project"})
		return
	}
}

func getAllProjectForFilter() ([]models.FilterNeededData, error) {
	var allData []models.FilterNeededData

	data, err := config.DB.Query("SELECT id, title, description, skills FROM project")

	if err != nil {
		return []models.FilterNeededData{}, errors.New("Server is unable to execute query to the database")
	}

	for data.Next() {
		var dbData models.FilterNeededData
		if err := data.Scan(&dbData.ID, &dbData.Title, &dbData.Description, &dbData.Skill); err != nil {
			return []models.FilterNeededData{}, errors.New("Something is wrong with the database data")
		}

		allData = append(allData, dbData)
	}

	return allData, nil
}

// SearchProject => Filter search project in SPIRITS
// SearchProject godoc
// @Summary Search and filter project here
// @Produce json
// @Tags Project
// @Param page query int64 true "page"
// @Param size query int64 true "size"
// @Param keyword query string false "Keyword"
// @Param sort query string false "Sort"
// @Param filter query string false "Filter Skills"
// @Success 200 {object} models.ResponseOKSearchProject
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /searchProject [get]
func SearchProject(c *gin.Context) {
	sortChoice := []string{"newest", "highestprice", "lowestprice"}
	interfaceChoice := make([]interface{}, len(sortChoice))
	for i, v := range sortChoice {
		interfaceChoice[i] = v
	}

	// var filteredID []string
	var wordFilter, skillFilter string
	var sortFilter string
	keyParam, okKey := c.Request.URL.Query()["keyword"]
	sortParam, okSort := c.Request.URL.Query()["sort"]
	filterParam, okFilter := c.Request.URL.Query()["filter"]

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

	if !okKey || len(keyParam[0]) < 1 {
		wordFilter = ""
	} else {
		wordFilter = keyParam[0]
	}
	if !okFilter || len(filterParam[0]) < 1 {
		skillFilter = ""
	} else {
		skillFilter = filterParam[0]
	}
	if !okSort || len(sortParam[0]) < 1 {
		sortFilter = ""
	} else {
		sortFilter = strings.ToLower(sortParam[0])
	}

	correctSortFilter := helpers.Contains(interfaceChoice, sortFilter)

	if sortFilter != "" {
		if !correctSortFilter {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Sorting name not recognized"})
			return
		}
	}

	allData, err := getAllProjectForFilter()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// filter project id based on title
	filteredID, err := filterData(allData, wordFilter, skillFilter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// get all the project where it has been filtered
	if len(filteredID) > 0 {
		stringID := strings.Trim(strings.Replace(fmt.Sprint(filteredID), " ", ",", -1), "[]")
		query, err := helpers.SettingInQueryWithID("project", stringID, "id, title, description, price")

		// get the project list that is listed
		query = query + " AND status=\"Listed\""

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

		switch sortFilter {
		case "newest":
			query = query + " ORDER BY id DESC"
		case "highestprice":
			query = query + " ORDER BY price DESC"
		case "lowestprice":
			query = query + " ORDER BY price ASC"
		default:
			query = query + " ORDER BY id DESC"
		}

		// include page and size
		var startingRecordNumber = page * size
		var endingRecordNumber = startingRecordNumber + size
		query = query + " LIMIT " + strconv.Itoa(startingRecordNumber) + "," + strconv.Itoa(endingRecordNumber)

		filteredProjectData, err := config.DB.Query(query)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server is unable to execute query to the database"})
			return
		}

		var project []models.SearchProjectQuery

		for filteredProjectData.Next() {
			var row models.SearchProjectQuery
			if err := filteredProjectData.Scan(&row.ID, &row.Title, &row.Description, &row.Price); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Something is wrong with the database data"})
				return
			}
			project = append(project, row)

		}

		if filteredProjectData.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Filter Successful",
			"data":    project})
	} else if sortFilter == "" && skillFilter == "" && wordFilter == "" {
		// this is when the search only page and size
		var startingRecordNumber = page * size
		var endingRecordNumber = startingRecordNumber + size

		result, err := config.DB.Query("SELECT id, title, description, price FROM project WHERE status=? ORDER BY ID DESC LIMIT ?,?", "Listed", startingRecordNumber, endingRecordNumber)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server is unable to execute query to the database"})
			return
		}
		var project []models.SearchProjectQuery

		for result.Next() {
			var row models.SearchProjectQuery
			if err := result.Scan(&row.ID, &row.Title, &row.Description, &row.Price); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Something is wrong with the database data"})
				return
			}
			project = append(project, row)
		}

		if result.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Filter Successful",
			"data":    project})
	} else if sortFilter != "" && wordFilter == "" {
		query := "SELECT id, title, description, price FROM project"

		switch sortFilter {
		case "newest":
			query = query + " ORDER BY id DESC"
		case "highestprice":
			query = query + " ORDER BY price DESC"
		case "lowestprice":
			query = query + " ORDER BY price ASC"
		}

		var startingRecordNumber = page * size
		var endingRecordNumber = startingRecordNumber + size
		query = query + " LIMIT " + strconv.Itoa(startingRecordNumber) + "," + strconv.Itoa(endingRecordNumber)

		filteredProjectData, err := config.DB.Query(query)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server is unable to execute query to the database"})
			return
		}

		var project []models.SearchProjectQuery

		for filteredProjectData.Next() {
			var row models.SearchProjectQuery
			if err := filteredProjectData.Scan(&row.ID, &row.Title, &row.Description, &row.Price); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Something is wrong with the database data"})
				return
			}
			project = append(project, row)

		}

		if filteredProjectData.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Filter Successful",
			"data":    project})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Filter Successful",
			"data":    []models.SearchProjectQuery{}})
	}
}

func filterData(data []models.FilterNeededData, keyword string, skill string) ([]int, error) {
	id := []int{}
	skillProjectID := []int{}

	if err != nil {
		return id, err
	}

	for i := 0; i < len(data); i++ {
		// title search based on keyword
		if keyword != "" && strings.Contains(strings.ToLower(data[i].Title), strings.ToLower(keyword)) {
			id = append(id, data[i].ID)
		}

		// description search based on keyword
		if keyword != "" && strings.Contains(strings.ToLower(data[i].Description), strings.ToLower(keyword)) {
			id = append(id, data[i].ID)
		}

		// skill search based on id filtered
		if skill != "" {
			arrSkill := helpers.SplitComma(data[i].Skill)
			arrFilteredSkill := helpers.SplitComma(skill)

			var t2 = []int{}
			for _, i := range arrFilteredSkill {
				j, err := strconv.Atoi(i)
				if err != nil {
					panic(err)
				}
				t2 = append(t2, j)
			}

			err = helpers.SkillList(t2)

			if err != nil {
				if err.Error() == "not exist" {
					return []int{}, errors.New("There is skill value that does not exist in the database id")
				}

				return []int{}, err
			}

			// find intersect value
			IntersectValue := helpers.FindDuplicateString(arrSkill, arrFilteredSkill)

			var integersIntersectValue = []int{}

			for _, i := range IntersectValue {
				j, err := strconv.Atoi(i)
				if err != nil {
					return []int{}, errors.New("Server is unable to convert string to integer")
				}
				integersIntersectValue = append(integersIntersectValue, j)
			}

			integersIntersectValue = helpers.RemoveDuplicateIntegerArray(integersIntersectValue)

			if len(integersIntersectValue) > 0 {
				skillProjectID = append(skillProjectID, data[i].ID)
			}
		}
	}

	// all filter condition
	if skill != "" && keyword != "" {
		id = helpers.FindDuplicateInteger(id, skillProjectID)
	}

	// if there is no keyword
	if keyword == "" && len(skillProjectID) > 0 {
		id = skillProjectID
	}

	id = helpers.RemoveDuplicateIntegerArray(id)

	return id, nil
}
