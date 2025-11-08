package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"unique;not null" json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

type URL struct {
	gorm.Model
	LongURL    string `gorm:"not null" json:"long_url"`
	ShortCode  string `gorm:"uniqueIndex;not null" json:"short_code"`
	UserID     *uint  `gorm:"not null" json:"user_id"`
	User       User   `gorm:"foreignKey:UserID" json:"user"`
	ClickCount int    `gorm:"default:0" json:"click_count"`
}
