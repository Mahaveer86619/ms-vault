package services

import "github.com/Mahaveer86619/ms/auth/pkg/models"

type AuthService struct {
	AvatarService *AvatarService
}

func NewAuthService(avatarService *AvatarService) *AuthService {
	return &AuthService{
		AvatarService: avatarService,
	}
}

func (s *AuthService) RegisterWithEmail(
	username string,
	email string,
	password string,
) (*models.UserProfile, error) {

	hash := s.AvatarService.GenerateHash(email)

	avatarURL := s.AvatarService.GetAvatarURL(hash)

	user := models.UserProfile{
		Username:  username,
		Email:     email,
		Password:  password,
		AvatarUrl: avatarURL,
	}

	// Placeholder for GORM create operation (assuming db.DB is initialized GORM)
	// if err := db.DB.Create(&user).Error; err != nil {
	// 	return nil, err
	// }

	// NOTE: You'll need to update db.DB from `*sql.DB` to a GORM connection (`*gorm.DB`)
	// or use your current connection setup for this to work with your models.UserProfile.

	return &user, nil
}

func (s *AuthService) LoginWithEmail() {}
