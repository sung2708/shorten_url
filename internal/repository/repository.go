package repository

import "github.com/sung2708/shorten_url/internal/model"

type UrlRepository interface {
	Save(url *model.URL) (*model.URL, error)
	Find(code string) (*model.URL, error)
	Delete(code string) error
	//	Update(oldCode string, newURL *string, newCode *string) (*model.URL, error)
	FindByUserID(userID uint) ([]*model.URL, error)
}

type UserRepository interface {
	Save(u *model.User) error
	Update(u *model.User) error
	Delete(u *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
}

type OTPRepository interface {
	Save(otp *model.OTP) error
	Find(userID uint, code string) (*model.OTP, error)
	Delete(userID uint, code string) error
	DeleteExpired() error
	CanResend(userID uint) (bool, error) // Kiểm tra gửi lại OTP
	UpdateLastSent(userID uint) error
	IncrementAttempts(userID uint, code string) error
}
