package repository

import (
	"github.com/sung2708/shorten_url/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
func (ur *UserRepository) Save(u *model.User) error {
	return ur.db.Create(u).Error
}

func (ur *UserRepository) Update(u *model.User) error {
	return ur.db.Save(u).Error
}

func (ur *UserRepository) Delete(u *model.User) error {
	return ur.db.Delete(u).Error
}

func (ur *UserRepository) FindById(id uint) (*model.User, error) {
	var u model.User
	return &u, ur.db.Where("id = ?", id).First(&u).Error
}
