package auth

import "marketplace-app/repositories"

type AuthHandler struct {
	repository *repositories.UserRepository
}

func NewHandler(userRepo *repositories.UserRepository) *AuthHandler {
	return &AuthHandler{repository: userRepo}
}
