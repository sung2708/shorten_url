package repository

import (
	"github.com/sung2708/shorten_url/internal/model"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (ur *UserRepositoryImpl) Save(u *model.User) error {
	return ur.db.Create(u).Error
}

func (ur *UserRepositoryImpl) Update(u *model.User) error {
	return ur.db.Save(u).Error
}

func (ur *UserRepositoryImpl) Delete(u *model.User) error {
	return ur.db.Delete(u).Error
}

func (ur *UserRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	var u model.User
	err := ur.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
