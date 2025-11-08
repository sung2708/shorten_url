package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/sung2708/shorten_url/internal/model"
	"github.com/sung2708/shorten_url/internal/repository"
)

type UrlServiceImpl struct {
	repo repository.UrlRepository
}

func NewUrlService(repo repository.UrlRepository) *UrlServiceImpl {
	return &UrlServiceImpl{repo: repo}
}

func (u *UrlServiceImpl) Shorten(url string, userID *uint) (*model.URL, error) {
	h1 := sha1.New()
	h1.Write([]byte(url))
	code := hex.EncodeToString(h1.Sum(nil))[:6]
	normalizedURL := url
	if !strings.HasPrefix(normalizedURL, "http://") && !strings.HasPrefix(normalizedURL, "https://") {
		normalizedURL = "https://" + normalizedURL
	}

	newURL := &model.URL{
		LongURL:   url,
		ShortCode: code,
		UserID:    userID,
	}
	createURL, err := u.repo.Save(newURL)
	if err != nil {
		return nil, err
	}
	return createURL, nil
}

func (u *UrlServiceImpl) GetById(code string) (*model.URL, error) {
	url, err := u.repo.Find(code)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (u *UrlServiceImpl) FindByUserID(userID uint) ([]*model.URL, error) {
	return u.repo.FindByUserID(userID)
}

func (u *UrlServiceImpl) DeleteLink(code string, userID uint) error {
	url, err := u.repo.Find(code)
	if err != nil {
		return errors.New("url not found")
	}
	if url.UserID == nil || *url.UserID != userID {
		return errors.New("user is not own link")
	}
	return u.repo.Delete(code)
}
