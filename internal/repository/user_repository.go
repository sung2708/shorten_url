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

func (r *UserRepositoryImpl) Save(u *model.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepositoryImpl) Update(u *model.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepositoryImpl) Delete(u *model.User) error {
	return r.db.Delete(u).Error
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepositoryImpl) FindByID(id uint) (*model.User, error) {
	var u model.User
	err := r.db.First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
