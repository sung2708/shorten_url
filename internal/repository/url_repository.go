package repository

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sung2708/shorten_url/internal/model"
	"gorm.io/gorm"
)

type URLRepository struct {
	db  *gorm.DB
	rdb *redis.Client
	ctx context.Context
}

func NewURLRepository(db *gorm.DB, rdb *redis.Client) *URLRepository {
	return &URLRepository{
		db:  db,
		rdb: rdb,
		ctx: context.Background(),
	}
}
func (r *URLRepository) Save(url *model.URL) error {
	if err := r.db.Create(url).Error; err != nil {
		return err
	}
	err := r.rdb.Set(r.ctx, url.ShortCode, url.LongURl, 24*time.Hour).Err()
	if err != nil {
		log.Println("Could not save URL: ", url.ShortCode, err)
	}
	return nil
}

func (r *URLRepository) Find(code string) (*model.URL, error) {
	longURL, err := r.rdb.Get(r.ctx, code).Result()
	if err == nil {
		log.Println("Catch hit", code)
		return &model.URL{
			ShortCode: code,
			LongURl:   longURL,
		}, nil
	}
	log.Println("Catch miss", code)
	var url model.URL
	result := r.db.Where("code = ?", code).First(&url).Error

	if result.Error != nil {
		log.Println("Could not find URL: ", code)
	}
	err = r.rdb.Set(r.ctx, url.ShortCode, url.LongURl, 24*time.Hour).Err()
	if err != nil {
		log.Println("Could not save URL: ", url.ShortCode, err)
	}
	return &url, nil
}
