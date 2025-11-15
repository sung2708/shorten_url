package repository

import (
	"errors"
	"time"

	"github.com/sung2708/shorten_url/internal/model"
	"gorm.io/gorm"
)

type OTPRepositoryImpl struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &OTPRepositoryImpl{db: db}
}

func (r *OTPRepositoryImpl) Save(otp *model.OTP) error {
	return r.db.Create(otp).Error
}

func (r *OTPRepositoryImpl) Find(userID uint, code string) (*model.OTP, error) {
	var otp model.OTP
	err := r.db.Where("user_id = ? AND code = ?", userID, code).First(&otp).Error
	if err != nil {
		return nil, err
	}
	if time.Now().After(otp.ExpiresAt) {
		return nil, gorm.ErrRecordNotFound
	}
	return &otp, nil
}

func (r *OTPRepositoryImpl) Delete(userID uint, code string) error {
	return r.db.Where("user_id = ? AND code = ?", userID, code).Delete(&model.OTP{}).Error
}

func (r *OTPRepositoryImpl) DeleteExpired() error {
	return r.db.Where("expires_at <= ?", time.Now()).Delete(&model.OTP{}).Error
}

func (r *OTPRepositoryImpl) CanResend(userID uint) (bool, error) {
	var otp model.OTP
	err := r.db.Where("user_id = ?", userID).Order("last_sent DESC").First(&otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}

	if time.Since(otp.LastSent) < time.Minute {
		return false, nil
	}
	return true, nil
}

func (r *OTPRepositoryImpl) UpdateLastSent(userID uint) error {
	var otp model.OTP
	err := r.db.Where("user_id = ?", userID).Order("last_sent DESC").First(&otp).Error
	if err != nil {
		return err
	}
	return r.db.Model(&otp).Update("last_sent", time.Now()).Error
}

func (r *OTPRepositoryImpl) IncrementAttempts(userID uint, code string) error {
	return r.db.Model(&model.OTP{}).
		Where("user_id = ? AND code = ?", userID, code).
		Update("attempts", gorm.Expr("attempts + ?", 1)).Error
}
