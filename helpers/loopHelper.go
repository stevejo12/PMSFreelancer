package helpers

// Contains => helper to find if string value is in the array
// arr => array for checking
// s => the string you want to check
func Contains(arr []string, s interface{}) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}
