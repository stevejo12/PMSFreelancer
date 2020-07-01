package models

// GetTransactionMutationRequest func
type GetTransactionMutationRequest struct {
	ID            string  `json:"id"`
	BankID        string  `json:"bank_id"`
	AccountNumber string  `json:"account_number"`
	BankType      string  `json:"bank_type"`
	Date          string  `json:"date"`
	Amount        int     `json:"amount"`
	Description   string  `json:"description"`
	Type          string  `json:"type"`
	Balance       float64 `json:"balance"`
}

type WithdrawListData struct {
	ID            int     `json:"id"`
	Amount        float64 `json:"amount"`
	Name          string  `json:"name"`
	AccountNumber int     `json:"accountNumber"`
}

type AllWithdrawListData struct {
	WithdrawListData
	Username string `json:"username"`
}

type WithdrawParameter struct {
	AccountName   string  `json:"accountName"`
	AccountNumber int     `json:"accountNumber"`
	Amount        float64 `json:"amount"`
}

type DepositParameter struct {
	AccountName   string  `json:"accountName"`
	AccountNumber int     `json:"accountNumber"`
	Amount        float64 `json:"amount"`
}

type ResponseDeposit struct {
	Amount          int    `json:"amount"`
	TransferAccount string `json:"transferAccount"`
}
