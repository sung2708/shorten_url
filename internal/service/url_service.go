package service

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/sung2708/shorten_url/internal/model"
	"github.com/sung2708/shorten_url/internal/repository"
)

type UrlServiceImpl struct {
	repo *repository.URLRepository
}

func NewUrlService(repo *repository.URLRepository) *UrlServiceImpl {
	return &UrlServiceImpl{repo: repo}
}

func (u *UrlServiceImpl) Shorten(url string, userID *uint) (*model.URL, error) {
	h1 := sha1.New()
	h1.Write([]byte(url))
	code := hex.EncodeToString(h1.Sum(nil))[:6]

	newURL := &model.URL{
		LongURl:   url,
		ShortCode: code,
		UserID:    userID,
	}
	if err := u.repo.Save(newURL); err != nil {
		return nil, err
	}
	return newURL, nil
}

func (u *UrlServiceImpl) GetById(code string) (*model.URL, error) {
	url, err := u.repo.Find(code)
	if err != nil {
		return nil, err
	}
	return url, nil
}
