package helpers

import (
	"errors"
	"strings"
)

// SplittingFullname => Separate Fullname into first name and last name
func SplittingFullname(fn string) (string, string, error) {
	splitName := SplitSpace(fn)

	if len(splitName) <= 0 {
		return "", "", errors.New("empty")
	} else if len(splitName) == 1 {
		return splitName[0], splitName[0], nil
	}

	firstname := strings.Join(splitName[:len(splitName)-1], " ")
	lastname := splitName[len(splitName)-1]

	return firstname, lastname, nil
}
