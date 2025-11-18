package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSecurity struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex;not null"`

	IsLocked            bool `gorm:"default:false"`
	LockedUntil         time.Time
	IsSuspended         bool `gorm:"default:false"`
	FailedLoginAttempts int  `gorm:"default:0"`

	TokenVersion int `gorm:"default:1"`
}
