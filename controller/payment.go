package controller

import (
	// "PMSFreelancer/config"
	// "PMSFreelancer/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
)

var accNumber = "2440277481"

var mootaAPIToken = "aWC0mYOl1njZ5Owu2ItnMCgVty6OId6CFVCWrOs2z1eSTu31iy"

const baseAPIMoota = "https://app.moota.co/"

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
		fmt.Println("error depositing money")
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
// @Summary Deposit Money
// @Produce json
// @Accept  json
// @Tags Payment
// @Param token header string true "Token Header"
// @Param Data body models.DepositParameter true "Data Format to deposit money"
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
		defer data.Close()

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

func checkDepositHistory(data []models.GetTransactionMutationRequest) error {
	for i := 0; i < len(data); i++ {
		floatAmount, errAmount := strconv.ParseFloat(data[i].Amount, 64)
		if errAmount != nil {
			fmt.Println("data amount convertion error: ", errAmount.Error())
		}

		intAmount := int(floatAmount)
		amount := intAmount
		// date := data[i].Date

		data, err := config.DB.Query("SELECT id, user_id, amount FROM payment_request WHERE amount=? AND type=\"Deposit\" AND status=\"Pending\"", amount)
		defer data.Close()

		if err != nil {
			return errors.New("Server is unable to retrieve data from the database")
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

// GetUserWithdrawRequest => User Withdraw Request
// GetUserWithdrawRequest godoc
// @Summary All user pending withdraw requests
// @Produce json
// @Tags Payment
// @Param token header string true "Token Header"
// @Success 200 {object} models.ResponseOKAllWithdrawRequest
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /userWithdrawRequest [get]
func GetUserWithdrawRequest(c *gin.Context) {
	id := idToken

	query := "SELECT id, amount, name, account_number FROM payment_request WHERE type=\"Withdraw\" AND status=\"Pending\" AND user_id=" + id

	resp, err := config.DB.Query(query)
	defer resp.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to retrieve withdraw list"})
		return
	}

	allData := []models.WithdrawListData{}

	for resp.Next() {
		data := models.WithdrawListData{}
		if err := resp.Scan(&data.ID, &data.Amount, &data.Name, &data.AccountNumber); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		allData = append(allData, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully get user withdraw requests",
		"data":    allData})
}

// DeleteWithdrawRequest => Deleting User Education
// DeleteWithdrawRequest godoc
// @Summary Deleting User Education
// @Accept  json
// @Tags Payment
// @Param token header string true "Token Header"
// @Param id path int64 true "Withdraw Request ID"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /deleteUserWithdrawRequest/{id} [delete]
func DeleteWithdrawRequest(c *gin.Context) {
	id := idToken
	requestID := c.Param("id")

	var userID, requestType, status string
	var amount float64

	query := "SELECT amount, user_id, type, status FROM payment_request WHERE id=" + requestID

	err := config.DB.QueryRow(query).Scan(&amount, &userID, &requestType, &status)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to get withdraw request data"})
		return
	}

	// make sure the real user logged in that request the withdraw
	if userID != id {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "This is not the owner who request the withdraw"})
		return
	}

	if requestType != "Withdraw" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "This is not the withdraw request"})
		return
	}

	if status != "Pending" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "This withdraw has been done"})
		return
	}

	var balance float64
	queryBalance := "SELECT balance FROM login WHERE id=" + id
	err = config.DB.QueryRow(queryBalance).Scan(&balance)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User ID doesn't exist"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to retrieve user balance"})
		return
	}

	// return the money to user balance
	balance = balance + amount

	updateBalanceQuery := "UPDATE login SET balance=" + fmt.Sprintf("%f", balance) + " WHERE id=" + id

	_, err = config.DB.Exec(updateBalanceQuery)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to update user balance"})
		return
	}

	// delete withdraw request
	queryDelete := "DELETE FROM payment_request WHERE id=" + requestID
	_, err = config.DB.Exec(queryDelete)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to delete withdraw request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully remove withdraw request"})
}

// GetAllWithdrawRequest => List of Pending Withdraw Request
// GetAllWithdrawRequest godoc
// @Summary Admin: Getting all pending withdraw requests
// @Produce json
// @Tags Payment
// @Param token header string true "Token Header"
// @Success 200 {object} models.ResponseOKAllWithdrawRequest
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /allWithdrawRequests [get]
func GetAllWithdrawRequest(c *gin.Context) {
	query := "SELECT id, amount, name, account_number, user_id FROM payment_request WHERE type=\"Withdraw\" AND status=\"Pending\""

	resp, err := config.DB.Query(query)
	defer resp.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to retrieve withdraw list"})
		return
	}

	allData := []models.AllWithdrawListData{}

	for resp.Next() {
		data := models.AllWithdrawListData{}
		var userID int
		if err := resp.Scan(&data.ID, &data.Amount, &data.Name, &data.AccountNumber, &userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Something is wrong with the database data"})
			return
		}

		var username string
		err := config.DB.QueryRow("SELECT username FROM login WHERE id=?", userID).Scan(&username)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server can't find the user id who requested the withdraw"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Server can't find username of the user"})
			return
		}

		data.Username = username

		allData = append(allData, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully get all withdraw request",
		"data":    allData})
}

// CompleteWithdrawRequest => Complete Withdraw Request By ID
// CompleteWithdrawRequest godoc
// @Summary Admin: Complete Withdraw Request By ID
// @Produce json
// @Tags Payment
// @Param token header string true "Token Header"
// @Param id path int64 true "Withdraw Request ID"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /completeWithdrawRequest/{id} [put]
func CompleteWithdrawRequest(c *gin.Context) {
	requestID := c.Param("id")
	query := "SELECT type, status FROM payment_request WHERE id=" + requestID

	var reqType, status string
	err := config.DB.QueryRow(query).Scan(&reqType, &status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to retrieve request information"})
		return
	}

	if reqType != "Withdraw" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "This is not the withdraw request"})
		return
	}

	if status != "Pending" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "This withdraw has been done"})
		return
	}

	updateBalanceQuery := "UPDATE payment_request SET status=\"Done\" WHERE id=" + requestID

	_, err = config.DB.Exec(updateBalanceQuery)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to update payment request to done"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully complete withdraw request"})
}

// WithdrawMoney => User request to withdraw balance
// WithdrawMoney godoc
// @Summary Withdraw money
// @Produce json
// @Accept  json
// @Tags Payment
// @Param token header string true "Token Header"
// @Param Data body models.WithdrawParameter true "Data Format to withdraw money"
// @Success 200 {object} models.ResponseOKWithdrawMoney
// @Failure 400 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /submitWithdrawRequest [post]
func WithdrawMoney(c *gin.Context) {
	id := idToken

	param := models.WithdrawParameter{}

	err := c.BindJSON(&param)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	var balance float64
	getBalanceQuery := "SELECT balance FROM login WHERE id=" + id
	err = config.DB.QueryRow(getBalanceQuery).Scan(&balance)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "User ID doesn't exist"})
		return
	} else if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to retrieve user balance"})
		return
	}

	if balance < param.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Balance is not enough"})
		return
	}

	// deduct the balance
	balance = balance - param.Amount

	query := "UPDATE login SET balance=" + fmt.Sprintf("%f", balance) + " WHERE id=" + id

	_, err = config.DB.Exec(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to update user balance"})
		return
	}

	withdrawQuery := "INSERT INTO payment_request(amount, type, name, account_number, user_id) VALUES"
	withdrawQuery = withdrawQuery + "(" + fmt.Sprintf("%f", param.Amount) + ", \"Withdraw\", \"" + param.AccountName + "\", " + strconv.Itoa(param.AccountNumber) + ", " + id + ")"

	_, err = config.DB.Exec(withdrawQuery)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to create withdraw request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully submit withdraw request"})
}
