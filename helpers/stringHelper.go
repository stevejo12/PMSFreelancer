package helpers

import (
	"fmt"
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
