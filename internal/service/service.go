package service

import "github.com/sung2708/shorten_url/internal/model"

type UrlService interface {
	Shorten(url string, userID *uint) (*model.URL, error)
	GetById(shortCode string) (*model.URL, error)
	FindByUserID(userID uint) ([]*model.URL, error)
	DeleteLink(code string, userID uint) error
}

type UserService interface {
	Register(u model.User) (*model.User, error)
	Login(email string, password string) (string, error)
}
