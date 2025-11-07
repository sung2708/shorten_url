package service

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sung2708/shorten_url/internal/model"
	"github.com/sung2708/shorten_url/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserServiceImpl {
	return &UserServiceImpl{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserServiceImpl) Register(u model.User) (*model.User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	u.Password = string(hashPassword)
	if err := s.repo.Save(&u); err != nil {
		log.Println(err)
		return nil, err
	}
	u.Password = ""
	return &u, nil
}

func (s *UserServiceImpl) Login(email string, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("Wrong email or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("Wrong email or password")
	}

	tokenString, err := s.generateJWT(user.ID)

	if err != nil {
		return "", err
	}

	user.Password = ""
	return tokenString, nil
}

func (s *UserServiceImpl) generateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
