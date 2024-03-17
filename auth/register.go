package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"marketplace-app/entities"
	"marketplace-app/utils"
	"net/http"
	"strings"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=5,max=15"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type RegisterResponse struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request data: %s", err.Error())})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
		return
	}

	hashedPw, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUser := &entities.User{
		Model:    entities.Model{ID: uuid.New()},
		Username: strings.ToLower(req.Username),
		Name:     strings.ToLower(req.Name),
		Password: hashedPw,
	}

	if err := h.repository.Create(newUser); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create user: %s", err.Error())})
		return
	}

	accessToken, err := utils.GenerateToken(newUser.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate access token: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data": RegisterResponse{
			Username:    newUser.Username,
			Name:        newUser.Name,
			AccessToken: accessToken,
		},
	})
}
