package controller

import (
	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
)

type testResponse struct {
	Name    string
	Email   string
	Address string
	City    string
	Join_at string
}

var accNumber = "2440277481"

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
	respData, err := callbackTransaction(req)

	if err != nil {
		fmt.Println("callbackerror")
		fmt.Println(err)
	}

	// tambahin buat add deposit nya
	err = checkDepositHistory(respData)

	if err != nil {
		fmt.Println("error depositting money")
		fmt.Println(err.Error())
	}

	// return success get
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully get mutation",
		"data":    respData})
}

func callbackTransaction(req string) ([]models.GetTransactionMutationRequest, error) {
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

// DepositMoney => User will request to top up balance
// DepositMoney godoc
// @Summary Adding User Project
// @Produce json
// @Accept  json
// @Tags Payment
// @Param token header string true "Token Header"
// @Param Data body models.CreateProject true "Data Format to add project"
// @Success 200 {object} models.ResponseOKDepositMoney
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /depositMoney [post]
func DepositMoney(c *gin.Context) {
	id := idToken

	param := models.DepositParameter{}

	err := c.BindJSON(&param)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	initialAmount := int(param.Amount)

	num := rand.Intn(1000)

	amountPaid := initialAmount + num

	// data check if there is a transaction with same unique number that has not been completed
	notUnique := true
	for notUnique {
		data, err := config.DB.Query("SELECT * FROM payment_request WHERE amount=? && type=\"Deposit\" && status=\"Pending\"", amountPaid)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server is unable to delete the data in the database"})
			return
		}

		if data.Next() {
			num = rand.Intn(1000)
			amountPaid = initialAmount + num
		} else {
			notUnique = false
		}
	}

	query := "INSERT INTO payment_request(amount, type, name, account_number, user_id) VALUES"
	query = query + "(" + strconv.Itoa(amountPaid) + ", \"Deposit\", \"" + param.AccountName + "\", " + strconv.Itoa(param.AccountNumber) + ", " + id + ")"

	_, err = config.DB.Exec(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})
		return
	}

	// response data
	resp := models.ResponseDeposit{}
	resp.Amount = amountPaid
	resp.TransferAccount = accNumber

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully Added Review",
		"data":    resp})
}

func TestFunction(c *gin.Context) {
	singleData := models.GetTransactionMutationRequest{}
	data := []models.GetTransactionMutationRequest{}

	singleData.ID = ""
	singleData.BankID = ""
	singleData.AccountNumber = ""
	singleData.BankType = "bca"
	singleData.Date = "30/06/2020"
	singleData.Amount = 20887
	singleData.Description = "TRSF E-BANKING"
	singleData.Type = "CR"
	singleData.Balance = 1900000

	data = append(data, singleData)
	err := checkDepositHistory(data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})
		return
	}
}

func checkDepositHistory(data []models.GetTransactionMutationRequest) error {
	for i := 0; i < len(data); i++ {
		amount := data[i].Amount
		// date := data[i].Date

		data, err := config.DB.Query("SELECT id, user_id, amount FROM payment_request WHERE amount=? AND type=\"Deposit\" AND status=\"Pending\"", amount)

		if err != nil {
			return errors.New("Server is unable to delete the data in the database")
		}

		var dataID, dataUserID, dataAmount int

		for data.Next() {
			if err := data.Scan(&dataID, &dataUserID, &dataAmount); err != nil {
				return errors.New("Something is wrong with the database data")
			}

			err = updateUserBalance(dataID, dataUserID, dataAmount)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func updateUserBalance(id int, userID int, amount int) error {
	// get userCurrent Balance
	var floatBalance float64

	err := config.DB.QueryRow("SELECT balance FROM login WHERE id=?", userID).Scan(&floatBalance)

	balance := int(floatBalance)

	if err != nil {
		return errors.New("Unable to retrieve user balance")
	}

	totalBalance := balance + amount

	query := "UPDATE login SET balance=" + strconv.Itoa(totalBalance)
	query = query + " WHERE id=" + strconv.Itoa(userID)

	_, err = config.DB.Exec(query)

	if err != nil {
		return errors.New("Unable to add user balance")
	}

	queryPayment := "UPDATE payment_request SET status=\"" + "Done" + "\""
	queryPayment = queryPayment + " WHERE id=" + strconv.Itoa(id)

	_, err = config.DB.Exec(queryPayment)

	if err != nil {
		return errors.New("Unable update status to done")
	}

	return nil
}
