package helpers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
)

// IsThisIDProjectOwner => Check if the user id is the project owner
func IsThisIDProjectOwner(projectID string, userID int) (bool, error) {
	var dbOwnerID int
	err := config.DB.QueryRow("SELECT owner_id FROM project WHERE id=?", projectID).Scan(&dbOwnerID)

	if err != nil {
		return false, errors.New("Server is unable to execute query to database")
	} else if dbOwnerID != userID {
		return false, nil
	}

	return true, nil
}

func IsThisMemberRegistered(projectID string, userID int) (bool, error) {
	var listInterested sql.NullString
	err := config.DB.QueryRow("SELECT interested_members FROM project WHERE id=?", projectID).Scan(&listInterested)

	if err != nil {
		return false, errors.New("Server is unable to execute query to database")
	}

	if listInterested.Valid {
		arrInterestedUser := SplitComma(listInterested.String)
		interfaceDataTypeInterestedUser := make([]interface{}, len(arrInterestedUser))
		for i := 0; i < len(interfaceDataTypeInterestedUser); i++ {
			interfaceDataTypeInterestedUser[i], err = strconv.Atoi(arrInterestedUser[i])
		}

		if len(arrInterestedUser) > 0 {
			exist := Contains(interfaceDataTypeInterestedUser, userID)

			if exist {
				return true, nil
			} else {
				return false, nil
			}
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}

func IsThisTheAcceptedMember(projectID string, userID int) (bool, error) {
	var dbMemberID int
	err := config.DB.QueryRow("SELECT accepted_memberid FROM project WHERE id=?", projectID).Scan(&dbMemberID)

	if err != nil {
		return false, errors.New("Server is unable to execute query to database")
	} else if dbMemberID != userID {
		return false, nil
	}

	return true, nil
}

// GetMemberList => List of interested freelancer for the project
func GetMemberList(id string) (string, error) {
	var dbListMember sql.NullString
	ok := config.DB.QueryRow("SELECT interested_members FROM project WHERE id=?", id).Scan(&dbListMember)

	if ok != nil {
		return "", errors.New("Server is unable to execute query to database")
	}

	if dbListMember.Valid {
		return dbListMember.String, nil
	} else {
		return "", nil
	}
}

func GetInterestedMemberNames(id string) ([]models.ProjectDetailInterestedMember, error) {
	imID, err := GetMemberList(id)
	allInterestedMember := []models.ProjectDetailInterestedMember{}

	query, err := SettingInQueryWithID("login", imID, "id, first_name, last_name")

	if err != nil {
		return []models.ProjectDetailInterestedMember{}, errors.New("Server has issues generating query")
	}

	data, err := config.DB.Query(query)

	if err != nil {
		return []models.ProjectDetailInterestedMember{}, errors.New("Server has issues executing query to the database")
	}

	for data.Next() {
		interestedMember := models.ProjectDetailInterestedMember{}
		var dbID, dbFirstName, dbLastName string
		if err := data.Scan(&dbID, &dbFirstName, &dbLastName); err != nil {
			return []models.ProjectDetailInterestedMember{}, errors.New("Something is wrong with the database data")
		}

		interestedMember.ID = dbID
		interestedMember.Fullname = dbFirstName + " " + dbLastName

		allInterestedMember = append(allInterestedMember, interestedMember)
	}

	return allInterestedMember, nil
}

func GetUserCompletedProject(id string) (int, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM project WHERE owner_id=?", id).Scan(&count)

	if err != nil {
		return -1, errors.New("Server unable to execute query to database")
	}

	return count, nil
}
