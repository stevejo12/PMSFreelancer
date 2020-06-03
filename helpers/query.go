package helpers

import "errors"

// SettingInQueryWithID => helping query problem with IN
// dbName => database name you want to select from
// param is the parameter. Format should be x,y,z
// this query is based on ids
func SettingInQueryWithID(dbName string, param string) (string, error) {
	arr := SplitComma(param)
	lengthArr := len(arr)
	initialQuery := "SELECT * FROM " + dbName + " WHERE id IN ("
	for i := 0; i < lengthArr; i++ {
		initialQuery = initialQuery + arr[i] + ","
	}

	var lengthString = len(initialQuery)

	if lengthArr > 0 && lengthString > 0 && initialQuery[lengthString-1] == ',' {
		initialQuery = initialQuery[:lengthString-1]
		initialQuery += ")"
	} else if lengthArr <= 0 {
		return "", errors.New("Parameter is empty")
	}

	return initialQuery, nil
}
