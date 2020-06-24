package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// IsEmpty => checking if string data is empty or not
func IsEmpty(data string) bool {
	if len(data) == 0 {
		return true
	}

	return false
}

func ConvertToString(val interface{}) string {
	str := fmt.Sprintf("%v", val)

	return str
}

// SplitComma => split string value separated by comma
func SplitComma(val string) []string {
	split := strings.Split(val, ",")

	if len(split) <= 0 {
		return []string{}
	} else {
		return split
	}
}

// SplitComma => split string value separated by space bar " "
func SplitSpace(val string) []string {
	split := strings.Split(val, " ")

	if len(split) <= 0 {
		return []string{}
	} else {
		return split
	}
}

// SplitDot => split string value separated by dot "."
func SplitDot(val string) []string {
	split := strings.Split(val, ".")

	if len(split) <= 0 {
		return []string{}
	} else {
		return split
	}
}

// SplitDot => split string value separated by dash "-"
func SplitDash(val string) []string {
	split := strings.Split(val, "-")

	if len(split) <= 0 {
		return []string{}
	} else {
		return split
	}
}

func ConvertStringToArrayInt(s string) ([]int, error) {
	var err error

	if s == "" {
		return []int{}, nil
	}

	arrInt := strings.Split(s, ",")
	arrIntSkill := make([]int, len(arrInt))

	fmt.Println(arrInt)
	fmt.Println(len(arrInt))

	for i := 0; i < len(arrInt); i++ {
		arrIntSkill[i], err = strconv.Atoi(arrInt[i])

		if err != nil {
			return []int{}, errors.New("Something wrong with convertion string to int")
		}
	}

	return arrIntSkill, nil
}
