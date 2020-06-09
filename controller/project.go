package controller

import (
	// "PMSFreelancer/config"
	// "PMSFreelancer/helpers"
	// "PMSFreelancer/models"
	"errors"
	"net/http"
	"strconv"
	"strings"

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
	query, err := helpers.SettingInQueryWithID("project_links", param, "*")

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

func AcceptProjectInterest(c *gin.Context) {
	id := c.Param("id")

	type acceptInterest struct {
		OwnerID      int
		FreelancerID int
	}

	var param acceptInterest

	err = c.Bind(&param)

	if err != nil || param.OwnerID == 0 || param.FreelancerID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
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
	owner, err := helpers.IsThisIDProjectOwner(id, param.OwnerID)

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
			// list the member as accepted member and update the status to "On Going"
			_, err = config.DB.Exec("UPDATE project SET status=?, accepted_memberid=? WHERE id=?", "On Going", param.FreelancerID, id)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Server is unable to execute query to the database"})
				return
			}

			// TO DO: Create Trello Board and store the link

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

func ReviewProject(c *gin.Context) {
	// project id
	id := c.Param("id")

	type reviewProject struct {
		UserID int
	}

	var param reviewProject

	err = c.Bind(&param)

	if err != nil || param.UserID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Data format is invalid"})
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

	freelancer, err := helpers.IsThisTheAcceptedMember(id, param.UserID)

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

func CompleteProject(c *gin.Context) {
	// project id
	id := c.Param("id")

	type reviewProject struct {
		OwnerID int
	}

	var param reviewProject

	err = c.Bind(&param)

	if err != nil || param.OwnerID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Data format is invalid"})
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

	owner, err := helpers.IsThisIDProjectOwner(id, param.OwnerID)

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

	type reviewProject struct {
		OwnerID int
	}

	var param reviewProject

	err = c.Bind(&param)

	if err != nil || param.OwnerID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Data format is invalid"})
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

	owner, err := helpers.IsThisIDProjectOwner(id, param.OwnerID)

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

func FilterProject(c *gin.Context) {
	var filteredID []string
	keyParam, ok := c.Request.URL.Query()["key"]

	if !ok || len(keyParam[0]) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Request Url should have key for search in it"})
		return
	}

	keyword := keyParam[0]

	allData, err := getAllProjectForFilter()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// filter project id based on title
	filteredID, err = filterData(allData, keyword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}

	// get all the project where it has been filtered
	if len(filteredID) > 0 {
		query, err := helpers.SettingInQueryWithID("project", strings.Join(filteredID, ","), "id, title, description, price")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error()})
			return
		}

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

func filterData(data []models.FilterNeededData, keyword string) ([]string, error) {
	var id []string
	var skillID []string

	allSkills, err := getAllSkills()

	if err != nil {
		return id, err
	}

	// filter skill name that matches the keyword
	for i := 0; i < len(allSkills); i++ {
		if strings.Contains(strings.ToLower(allSkills[i].Name), strings.ToLower(keyword)) {
			skillID = append(skillID, allSkills[i].ID)
		}
	}

	for i := 0; i < len(data); i++ {
		if strings.Contains(strings.ToLower(data[i].Title), strings.ToLower(keyword)) {
			id = append(id, data[i].ID)
		}

		if strings.Contains(strings.ToLower(data[i].Description), strings.ToLower(keyword)) {
			id = append(id, data[i].ID)
		}

		arrSkill := helpers.SplitComma(data[i].Skill)

		var skillMap = make(map[string]bool)

		for _, ele := range skillID {
			skillMap[ele] = true
		}

		for _, name := range arrSkill {
			if skillMap[name] {
				id = append(id, data[i].ID)
				break
			}
		}
	}

	// remove any duplicate id in the array
	id = helpers.RemoveDuplicateValueArray(id)

	return id, nil
}
