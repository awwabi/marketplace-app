package entities

type BankAccount struct {
	Model
	UserID        string `json:"user_id"`
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
}
