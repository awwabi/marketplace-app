package bank_account

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type UpdateBankAccountRequest struct {
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
}

func (bh *BankAccountHandler) UpdateBankAccount(c *gin.Context) {
	var request UpdateBankAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bankAccountID := c.Param("id")
	if err := uuid.Validate(bankAccountID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bank account ID"})
		return
	}

	bankAccount, err := bh.repository.FindByID(bankAccountID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bank account not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	if bankAccount.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this bank account"})
		return
	}

	bankAccount.BankName = request.BankName
	bankAccount.AccountName = request.BankAccountName
	bankAccount.AccountNumber = request.BankAccountNumber

	if err := bh.repository.Update(bankAccount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bankAccount)
}
