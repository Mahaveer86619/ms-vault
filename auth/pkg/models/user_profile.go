package models

import "gorm.io/gorm"

type UserProfile struct {
	gorm.Model

	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	Email     string `gorm:"unique;not null"`
	AvatarUrl string
	
}
