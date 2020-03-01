package helpers

import "fmt"

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
