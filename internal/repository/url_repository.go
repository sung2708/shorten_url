package repository

import (
	"github.com/sung2708/shorten_url/internal/model"
	"gorm.io/gorm"
)

type URLRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{db: db}
}
func (r *URLRepository) Save(url *model.URL) error {
	return r.db.Create(url).Error
}

func (r *URLRepository) Find(code string) (*model.URL, error) {
	var url model.URL
	err := r.db.Where("code = ?", code).First(&url).Error

	if err != nil {
		return nil, err
	}
	return &url, nil
}
