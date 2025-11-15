package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sung2708/shorten_url/internal/model"
	"github.com/sung2708/shorten_url/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repo                repository.UserRepository
	otpRepo             repository.OTPRepository
	notificationService NotificationService
	jwtSecret           string
}

func NewUserService(repo repository.UserRepository, otpRepo repository.OTPRepository, notificationService NotificationService, jwtSecret string) *UserServiceImpl {
	return &UserServiceImpl{
		repo:                repo,
		otpRepo:             otpRepo,
		notificationService: notificationService,
		jwtSecret:           jwtSecret,
	}
}

func (s *UserServiceImpl) Register(u model.User) (string, *model.User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}
	u.Password = string(hashPassword)

	if err := s.repo.Save(&u); err != nil {
		return "", nil, err
	}

	// tạo OTP và gửi email
	otp, err := s.notificationService.GenerateOTP(u.ID)
	if err == nil {
		// Ignore email sending errors during registration
		_ = s.notificationService.SendVerificationOTP(&u, otp)
	}

	tokenString, err := s.generateJWT(u.ID, false)
	if err != nil {
		return "", nil, err
	}
	u.Password = ""
	return tokenString, &u, nil
}

func (s *UserServiceImpl) Login(email string, password string) (string, *model.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("wrong email")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", nil, errors.New("wrong password")
	}

	// nếu chưa verify → gửi lại OTP & trả token is_verified=false
	if !user.IsActive {
		otp, err := s.notificationService.GenerateOTP(user.ID)
		if err == nil {
			// Ignore email sending errors during login
			_ = s.notificationService.SendVerificationOTP(user, otp)
		}
		token, err := s.generateJWT(user.ID, false)
		if err != nil {
			return "", nil, err
		}
		user.Password = ""
		return token, user, nil
	}

	// verified → cấp token full quyền
	token, err := s.generateJWT(user.ID, true)
	if err != nil {
		return "", nil, err
	}
	user.Password = ""
	return token, user, nil
}

func (s *UserServiceImpl) VerifyAccount(userID uint, otp string) (string, *model.User, error) {
	// VerifyOTP already activates the user if not active
	user, err := s.notificationService.VerifyOTP(userID, otp)
	if err != nil {
		return "", nil, err
	}

	// cấp token mới với quyền đầy đủ
	token, err := s.generateJWT(user.ID, true)
	if err != nil {
		return "", nil, err
	}
	user.Password = ""
	return token, user, nil
}

func (s *UserServiceImpl) generateJWT(userID uint, isVerified bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     userID,
		"is_verified": isVerified,
		"exp":         time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":         time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
