package bank_account

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (bh *BankAccountHandler) DeleteBankAccount(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bank account ID"})
		return
	}

	bankAccount, err := bh.repository.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this bank account"})
		return
	}

	err = bh.repository.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bank account deleted successfully"})
}
