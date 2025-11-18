package services

import (
	"errors"
	"time"

	"github.com/Mahaveer86619/ms/auth/pkg/db"
	"github.com/Mahaveer86619/ms/auth/pkg/models"
	"github.com/Mahaveer86619/ms/auth/pkg/utils"
	"github.com/Mahaveer86619/ms/auth/pkg/views"
	"gorm.io/gorm"
)

const (
	MaxLoginAttempts = 5
	LockDuration     = 15 * time.Minute
)

type AuthService struct {
	AvatarService *AvatarService
}

func NewAuthService(avatarService *AvatarService) *AuthService {
	return &AuthService{
		AvatarService: avatarService,
	}
}

func (s *AuthService) RegisterUser(username, email, password string) (*views.AuthResponse, error) {
	var existingUser models.UserProfile
	err := db.DB.Where("username = ?", username).Or("email = ?", email).First(&existingUser).Error

	if err == nil {
		return nil, errors.New("username already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	var user models.UserProfile
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		avatarHash := s.AvatarService.GenerateHash(email)
		user = models.UserProfile{
			Username:  username,
			Email:     email,
			Password:  hashedPwd,
			AvatarUrl: s.AvatarService.GetAvatarURL(avatarHash),
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		security := models.UserSecurity{UserID: user.ID}
		if err := tx.Create(&security).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) LoginWithUsername(username, password string) (*views.AuthResponse, error) {
	var user models.UserProfile
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	var security models.UserSecurity
	if err := db.DB.Where("user_id = ?", user.ID).First(&security).Error; err != nil {
		return nil, errors.New("system error: security profile missing")
	}

	if security.IsSuspended {
		return nil, errors.New("account is suspended")
	}

	if security.IsLocked {
		if time.Now().Before(security.LockedUntil) {
			return nil, errors.New("account is temporarily locked due to too many failed attempts")
		}
		s.resetLock(&security)
	}

	if err := utils.CheckPassword(user.Password, password); err != nil {
		s.handleFailedAttempt(&security)
		return nil, errors.New("invalid credentials")
	}

	if security.FailedLoginAttempts > 0 {
		s.resetLock(&security)
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*views.AuthResponse, error) {
	claims, err := utils.ValidateToken(refreshToken, "refresh")
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var user models.UserProfile
	if err := db.DB.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) handleFailedAttempt(sec *models.UserSecurity) {
	sec.FailedLoginAttempts++
	if sec.FailedLoginAttempts >= MaxLoginAttempts {
		sec.IsLocked = true
		sec.LockedUntil = time.Now().Add(LockDuration)
	}
	db.DB.Save(sec)
}

func (s *AuthService) resetLock(sec *models.UserSecurity) {
	sec.IsLocked = false
	sec.FailedLoginAttempts = 0
	sec.LockedUntil = time.Time{} // Zero time
	db.DB.Save(sec)
}

func (s *AuthService) generateAuthResponse(user models.UserProfile) (*views.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}
	return views.NewAuthResponse(user, accessToken, refreshToken), nil
}
