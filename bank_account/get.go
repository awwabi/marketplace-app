package bank_account

import (
	"github.com/gin-gonic/gin"
	"marketplace-app/entities"
	"net/http"
)

type BankAccountResponse struct {
	BankAccountID     string `json:"bankAccountId"`
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
}

func (bh *BankAccountHandler) GetBankAccount(c *gin.Context) {
	bankAccounts, err := bh.repository.FindByUserID(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapToResponse(bankAccounts))
}

func mapToResponse(bankAccounts []entities.BankAccount) []BankAccountResponse {
	var bankAccountResponses []BankAccountResponse
	for _, bankAccount := range bankAccounts {
		bankAccountResponses = append(bankAccountResponses, BankAccountResponse{
			BankAccountID:     bankAccount.ID.String(),
			BankName:          bankAccount.BankName,
			BankAccountName:   bankAccount.AccountName,
			BankAccountNumber: bankAccount.AccountNumber,
		})
	}
	return bankAccountResponses
}
