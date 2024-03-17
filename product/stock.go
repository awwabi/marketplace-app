package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"marketplace-app/utils"
	"net/http"
	"strings"
)

type UpdateStockRequest struct {
	Stock int `json:"stock" validate:"required,min=0"`
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	var req UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request data: %s", err.Error())})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
		return
	}

	productID := c.Param("id")
	product, err := h.repository.FindByID(productID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to find product: %s", err.Error())})
		return
	}

	// check user permission
	userID := c.GetString("userID")
	if product.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this product"})
		return
	}

	product.Stock = req.Stock
	if err := h.repository.UpdateStock(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update product: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, product)
}
