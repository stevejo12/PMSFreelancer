package helpers

import (
	"github.com/stevejo12/PMSFreelancer/config"
	// "PMSFreelancer/config"
	"errors"
)

// Contains => helper to find if string value is in the array
// arr => array for checking
// s => the string you want to check
func Contains(arr []interface{}, s interface{}) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func FindDuplicateString(a []string, b []string) []string {
	var shortest, longest *[]string
	if len(a) < len(b) {
		shortest = &a
		longest = &b
	} else {
		shortest = &b
		longest = &a
	}
	// Turn the shortest slice into a map
	var m map[string]bool
	m = make(map[string]bool, len(*shortest))
	for _, s := range *shortest {
		m[s] = false
	}
	// Append values from the longest slice that don't exist in the map
	var diff []string
	for _, s := range *longest {
		if _, ok := m[s]; ok {
			diff = append(diff, s)
			continue
		}
	}

	return diff
}

func FindDuplicateInteger(a []int, b []int) []int {
	var shortest, longest *[]int
	if len(a) < len(b) {
		shortest = &a
		longest = &b
	} else {
		shortest = &b
		longest = &a
	}
	// Turn the shortest slice into a map
	var m map[int]bool
	m = make(map[int]bool, len(*shortest))
	for _, s := range *shortest {
		m[s] = false
	}
	// Append values from the longest slice that don't exist in the map
	var diff []int
	for _, s := range *longest {
		if _, ok := m[s]; ok {
			diff = append(diff, s)
			continue
		}
	}

	return diff
}

func IsThisAttachmentAlreadyExistInDatabase(link string) (int, error) {
	querySelect := "SELECT id FROM project_links WHERE project_link=\"" + link + "\""
	data, err := config.DB.Query(querySelect)
	defer data.Close()

	if err != nil {
		return -1, errors.New("Server is unable to execute query to the database")
	}

	for data.Next() {
		var dataID int
		if err := data.Scan(&dataID); err != nil {
			return -1, errors.New("Server is unable to execute query to the database")
		}
		return dataID, nil
	}

	return -1, nil
}

func RemoveAttachmentThatIsDeletedByUser(attachment []string, projectID string) error {
	var allLinks []string

	query := "SELECT project_link FROM project_links WHERE project_id=" + projectID

	allProjectAttachment, err := config.DB.Query(query)
	defer allProjectAttachment.Close()

	if err != nil {
		return errors.New("Server is unable to execute query to the database")
	}

	for allProjectAttachment.Next() {
		var dataLink string
		if err := allProjectAttachment.Scan(&dataLink); err != nil {
			return errors.New("Server is unable to retrieve database data")
		}

		allLinks = append(allLinks, dataLink)
	}

	interfaceChoice := make([]interface{}, len(attachment))
	for i, v := range attachment {
		interfaceChoice[i] = v
	}

	for i := 0; i < len(allLinks); i++ {
		stillExist := Contains(interfaceChoice, allLinks[i])

		if !stillExist {
			queryDelete := "DELETE FROM project_links WHERE project_link=\"" + allLinks[i] + "\""
			_, err = config.DB.Exec(queryDelete)

			if err != nil {
				return errors.New("Server is unable to execute query to the database")
			}
		}
	}

	return nil
}
