package repository

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sung2708/shorten_url/internal/model"
	"gorm.io/gorm"
)

type URLRepositoryImpl struct {
	db  *gorm.DB
	rdb *redis.Client
	ctx context.Context
}

func NewURLRepository(db *gorm.DB, rdb *redis.Client) *URLRepositoryImpl {
	return &URLRepositoryImpl{
		db:  db,
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (r *URLRepositoryImpl) Save(url *model.URL) (*model.URL, error) {
	if err := r.db.Create(url).Error; err != nil {
		return nil, err
	}

	err := r.rdb.Set(r.ctx, url.ShortCode, url.LongURL, 24*time.Hour).Err()
	if err != nil {
		log.Println("Could not save URL: ", url.ShortCode, err)
	}

	return url, nil
}

func (r *URLRepositoryImpl) Find(code string) (*model.URL, error) {
	longURL, err := r.rdb.Get(r.ctx, code).Result()
	if err == nil {
		log.Println("Cache HIT:", code)
		return &model.URL{
			ShortCode: code,
			LongURL:   longURL,
		}, nil
	}

	log.Println("Cache MISS:", code)
	var url model.URL

	err = r.db.Where("short_code = ?", code).First(&url).Error

	// Sửa 5: PHẢI return ngay nếu lỗi (ví dụ: not found)
	if err != nil {
		log.Println("Could not find URL in DB: ", code)
		return nil, err
	}

	errCache := r.rdb.Set(r.ctx, url.ShortCode, url.LongURL, 24*time.Hour).Err()
	if errCache != nil {
		log.Println("Could not save URL: ", url.ShortCode, err)
	}

	return &url, nil
}

func (r *URLRepositoryImpl) Delete(code string) error {

	if err := r.db.Where("short_code = ?", code).Delete(&model.URL{}).Error; err != nil {
		return err
	}

	if err := r.rdb.Del(r.ctx, code).Err(); err != nil {
		log.Println("Could not delete URL: ", code, err)
	}

	return nil
}

func (r *URLRepositoryImpl) FindByUserID(userID uint) ([]*model.URL, error) {
	var links []*model.URL

	err := r.db.Where("user_id = ?", userID).Order("id desc").Find(&links).Error
	if err != nil {
		return nil, err
	}

	return links, nil
}
