package service

import "github.com/sung2708/shorten_url/internal/model"

type UrlService interface {
	Shorten(url string, userID *uint) (*model.URL, error)
	GetById(shortCode string) (*model.URL, error)
	FindByUserID(userID uint) ([]*model.URL, error)
	DeleteLink(code string, userID uint) error
	//UpdateLink(oldCode string, newURL *string, newCode *string, userID uint) (*model.URL, error)
}

type UserService interface {
	Register(u model.User) (string, *model.User, error)
	Login(email string, password string) (string, *model.User, error)
	VerifyAccount(userID uint, otp string) (string, *model.User, error)
}

type NotificationService interface {
	GenerateOTP(userID uint) (string, error)
	SendVerificationOTP(user *model.User, code string) error
	VerifyOTP(userID uint, code string) (*model.User, error)
	ResendOTP(userID uint) (string, error)
	SendPasswordResetEmail(user *model.User) error
}
