package services

import (
	"errors"

	"github.com/Mahaveer86619/lumi/pkg/db"
	"github.com/Mahaveer86619/lumi/pkg/models"
	"github.com/Mahaveer86619/lumi/pkg/views"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us UserService) GetUserDetails(id uint) (*views.UserDetailsResponse, error) {
	// Get user profile
	var user models.UserProfile
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Get user security details
	var userSecurity models.UserSecurity
	if err := db.DB.Where("user_id = ?", id).First(&userSecurity).Error; err != nil {
		return nil, errors.New("system error: security profile missing")
	}

	response := views.NewUserDetailsResponse(user, userSecurity)
	return response, nil
}
