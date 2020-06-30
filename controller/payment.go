package controller

import (
	// "PMSFreelancer/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/models"
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

// GetTransactionMutation get payment detail
func GetTransactionMutation(c *gin.Context) {
	var req string
	if err = c.Bind(&req); err != nil {
		fmt.Println("error binding")
		fmt.Println(err)
	}

	// call service buat rapihin data nya
	respData, err := CallbackTransaction(req)

	if err != nil {
		fmt.Println("callbackerror")
		fmt.Println(err)
	}

	// tambahin buat add deposit nya

	// return success get
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully get mutation",
		"data":    respData})
}

// CallbackTransaction => adjust the data structure
func CallbackTransaction(req string) ([]models.GetTransactionMutationRequest, error) {
	var data []models.GetTransactionMutationRequest
	req = strings.Replace(req, "\\", "", -1)
	req = strings.Replace(req, "\"[", "[", -1)
	req = strings.Replace(req, "]\"", "]", -1)
	err := json.Unmarshal([]byte(req), &data)
	if err != nil {
		fmt.Println("ERROR UNMARSHAL", err)
		// return []models.GetTransactionMutationRequest{}, errors.New("Error unmarshal data")
	}

	fmt.Println("***TRANSACTION-REQ**", req)
	fmt.Println("***TRANSACTION-DATA**", data)

	// ********TRANSACTION-REQ******* [{"mutation_id":"XX","bank_type":"bca","date":"30/06/2020","amount":"100000","description":"TRSF E-BANKING","type":"CR","balance":1900000}]
	// ********TRANSACTION-DATA******* [{  0 bca 30/06/2020 0 TRSF E-BANKING CR 1900000}]

	return data, nil
}
