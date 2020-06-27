package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type testResponse struct {
	Name    string
	Email   string
	Address string
	City    string
	Join_at string
}

var mootaAPIToken = "aWC0mYOl1njZ5Owu2ItnMCgVty6OId6CFVCWrOs2z1eSTu31iy"

const baseAPIMoota = "https://app.moota.co/"

func GetMootaProfile(c *gin.Context) {
	url := baseAPIMoota + "api/v1/profile"

	fmt.Println(url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+mootaAPIToken)
	res, err := client.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Something is wrong with the request"})
		return
	}

	if res.StatusCode == http.StatusOK {
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		userInfo := testResponse{}
		jsonErr := json.Unmarshal(body, &userInfo)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		fmt.Println(userInfo.Name)
		fmt.Println(userInfo.Email)
		fmt.Println(userInfo.Address)
		fmt.Println(userInfo.City)
		fmt.Println(userInfo.Join_at)
	}

	fmt.Println(res)
}
