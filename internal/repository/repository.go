package repository

import "github.com/sung2708/shorten_url/internal/model"

type UrlRepository interface {
	Save(url *model.URL) (*model.URL, error)
	Find(code string) (*model.URL, error)
	Delete(code string) error
	FindByUserID(userID uint) ([]*model.URL, error)
}

type UserRepository interface {
	Save(u *model.User) error
	Update(u *model.User) error
	Delete(u *model.User) error
	FindByEmail(email string) (*model.User, error)
}
