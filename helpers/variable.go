package helpers

import (
	"reflect"
)

func IsEmptyData(object interface{}) bool {
	//First check normal definitions of empty
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	//Then see if it's a struct
	if reflect.ValueOf(object).Kind() == reflect.Struct {
		// and create an empty copy of the struct object to compare against
		structIterator := reflect.ValueOf(object)
		for i := 0; i < structIterator.NumField(); i++ {
			val := structIterator.Field(i).Interface()
			varType := structIterator.Type().Field(i).Type.String()

			// check type int
			if varType == "int" {
				if val == 0 {
					return true
				}
			}

			// Check if the field is zero-valued, meaning it won't be updated
			if reflect.DeepEqual(val, reflect.Zero(structIterator.Field(i).Type()).Interface()) {
				return true
			}
		}
	}
	return false
}

func IsYearConsistFourNumber(startYear int, endYear int) bool {
	msg1 := startYear >= 1000 && startYear <= 9999
	msg2 := endYear >= 1000 && endYear <= 9999

	if !msg1 || !msg2 {
		return false
	}

	return true
}
