package views

import (
	"github.com/Mahaveer86619/ms/auth/pkg/models"
	"github.com/Mahaveer86619/ms/auth/pkg/utils"
)

type AuthRegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginWithUsername struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	ID           utils.MaskedId `json:"id"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	AvatarUrl    string         `json:"avatar_url"`
	Token        string         `json:"token"`
	RefreshToken string         `json:"refresh_token"`
}

func NewAuthResponse(user models.UserProfile, accessToken, refreshToken string) *AuthResponse {
	return &AuthResponse{
		ID:           utils.Mask(user.ID),
		Username:     user.Username,
		Email:        user.Email,
		AvatarUrl:    user.AvatarUrl,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
}
