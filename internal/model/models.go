package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string
}

type URL struct {
	gorm.Model
	LongURl    string `gorm:"not null"`
	ShortCode  string `gorm:"uniqueIndex;not null"`
	UserID     *uint
	User       User `gorm:"foreignKey:UserID;AssociationForeignKey:ID;References:User"`
	ClickCount int  `gorm:"default:0"`
}
