package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	ImgURL   string `json:"img_url"`
	IsActive bool   `gorm:"default:false" json:"is_active"`
}

type URL struct {
	gorm.Model
	LongURL    string `gorm:"not null" json:"long_url"`
	ShortCode  string `gorm:"uniqueIndex;not null" json:"short_code"`
	UserID     *uint  `json:"user_id"`
	User       *User  `gorm:"foreignKey:UserID" json:"user"`
	ClickCount int    `gorm:"default:0" json:"click_count"`
}
type OTP struct {
	gorm.Model
	UserID    uint      `gorm:"index;not null"`
	Code      string    `gorm:"size:6;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Attempts  int       `gorm:"default:0"`
	LastSent  time.Time `gorm:"not null"`
}
