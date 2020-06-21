package helpers

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
