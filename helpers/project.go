package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	// "github.com/stevejo12/PMSFreelancer/config"
	"PMSFreelancer/config"
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
		fmt.Println(err)
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
