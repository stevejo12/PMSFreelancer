package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	// "PMSFreelancer/config"
	// "PMSFreelancer/models"

	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
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

	// if there is no interested member in the project
	if imID == "" {
		return []models.ProjectDetailInterestedMember{}, nil
	}

	query, err := SettingInQueryWithID("login", imID, "id, first_name, last_name, picture")

	if err != nil {
		return []models.ProjectDetailInterestedMember{}, errors.New("Server has issues generating query")
	}

	data, err := config.DB.Query(query)

	if err != nil {
		return []models.ProjectDetailInterestedMember{}, errors.New("Server has issues executing query to the database")
	}

	for data.Next() {
		interestedMember := models.ProjectDetailInterestedMember{}
		if err := data.Scan(&interestedMember.ID, &interestedMember.FirstName, &interestedMember.LastName, &interestedMember.Picture); err != nil {
			return []models.ProjectDetailInterestedMember{}, errors.New("Something is wrong with the database data")
		}

		allInterestedMember = append(allInterestedMember, interestedMember)
	}

	return allInterestedMember, nil
}

func GetUserCompletedProject(id string) (int, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM project WHERE owner_id=? && status=\"Done\"", id).Scan(&count)

	if err != nil {
		return -1, errors.New("Server unable to execute query to database")
	}

	return count, nil
}

func GetProjectCategory(id int) (string, error) {
	name := ""

	err := config.DB.QueryRow("SELECT name FROM project_category WHERE id=?", id).Scan(&name)

	if err == sql.ErrNoRows {
		return "", errors.New("Category ID is not registered")
	} else if err != nil {
		return "", errors.New("Server unable to get information from database")
	}

	return name, nil
}

// IsUserEffectiveBalanceEnough => Check User Effective Balance
func IsUserEffectiveBalanceEnough(id int, projectPrice float64) (bool, error) {
	var balance float64
	err := config.DB.QueryRow("SELECT balance FROM login WHERE id=?", id).Scan(&balance)

	if err != nil {
		return false, errors.New("Server is unable to execute query get balance")
	}

	if balance < projectPrice {
		return false, nil
	}

	return true, nil
}

// MoveBalanceToFreezeBalance => Move to freezebalance
// to ensure money is secured for freelancer to be transfered after the completion
func MoveBalanceToFreezeBalance(id int, projectPrice float64) error {
	var freezeBalance, balance float64
	err := config.DB.QueryRow("SELECT balance, freeze_balance FROM login WHERE id=?", id).Scan(&balance, &freezeBalance)

	if err != nil {
		return errors.New("Server is unable to execute query get balance")
	}

	balance = balance - projectPrice
	freezeBalance = freezeBalance + projectPrice

	query := "UPDATE login SET freeze_balance=" + fmt.Sprintf("%f", freezeBalance)
	query = query + ", balance=" + fmt.Sprintf("%f", balance)
	query = query + " WHERE id=" + strconv.Itoa(id)

	_, err = config.DB.Exec(query)

	if err != nil {
		return errors.New("Server is unable to update user freeze balance")
	}

	return nil
}

// UpdateUserBalanceAfterProject => move balance from owner to freelancer
func UpdateUserBalanceAfterProject(projectID string) error {
	var price float64
	var ownerID, freelancerID int
	err := config.DB.QueryRow("SELECT price, owner_id, accepted_memberid FROM project WHERE id=?", projectID).Scan(&price, &ownerID, &freelancerID)

	if err != nil {
		return errors.New("Server is unable to retrieve project price")
	}

	var ownerFreezeBalance, freelancerBalance float64
	err = config.DB.QueryRow("SELECT freeze_balance FROM login WHERE id=?", ownerID).Scan(&ownerFreezeBalance)

	if err != nil {
		return errors.New("Server is unable to retrieve user balance")
	}

	err = config.DB.QueryRow("SELECT balance FROM login WHERE id=?", freelancerID).Scan(&freelancerBalance)

	if err != nil {
		return errors.New("Server is unable to retrieve user balance")
	}

	ownerFreezeBalance = ownerFreezeBalance - price
	freelancerBalance = freelancerBalance + price

	ownerQuery := "UPDATE login SET freeze_balance=" + fmt.Sprintf("%f", ownerFreezeBalance) + " WHERE id=" + strconv.Itoa(ownerID)

	_, err = config.DB.Exec(ownerQuery)

	if err != nil {
		return errors.New("Server is unable to update owner balance")
	}

	freelancerQuery := "UPDATE login SET balance=" + fmt.Sprintf("%f", freelancerBalance) + " WHERE id=" + strconv.Itoa(freelancerID)

	_, err = config.DB.Exec(freelancerQuery)

	if err != nil {
		return errors.New("Server is unable to update freelancer balance")
	}

	return nil
}
