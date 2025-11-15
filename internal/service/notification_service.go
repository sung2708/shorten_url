package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sung2708/shorten_url/internal/model"
	"github.com/sung2708/shorten_url/internal/repository"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type NotificationServiceImpl struct {
	repo        repository.OTPRepository
	userRepo    repository.UserRepository
	smtpHost    string
	smtpPort    int
	smtpUser    string
	smtpPass    string
	sender      string
	otpExpiry   time.Duration
	resendDelay time.Duration
}

func NewNotificationService(
	repo repository.OTPRepository,
	userRepo repository.UserRepository,
	smtpHost, smtpUser, smtpPass, sender string,
	smtpPort int,
) *NotificationServiceImpl {
	return &NotificationServiceImpl{
		repo:        repo,
		userRepo:    userRepo,
		smtpHost:    smtpHost,
		smtpPort:    smtpPort,
		smtpUser:    smtpUser,
		smtpPass:    smtpPass,
		sender:      sender,
		otpExpiry:   5 * time.Minute,
		resendDelay: 1 * time.Minute,
	}
}

func (s *NotificationServiceImpl) GenerateOTP(userID uint) (string, error) {
	code := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	otp := &model.OTP{
		UserID:    userID,
		Code:      code,
		ExpiresAt: time.Now().Add(s.otpExpiry),
		LastSent:  time.Now(),
	}
	if err := s.repo.Save(otp); err != nil {
		return "", err
	}
	return code, nil
}

func (s *NotificationServiceImpl) SendVerificationOTP(user *model.User, code string) error {
	subject := "Verify your account"
	body := fmt.Sprintf(`
	<p>Hi %s,</p>
	<p>Your verification code is: <b>%s</b></p>
	<p>This code expires in 5 minutes.</p>
	`, user.Email, code)
	return s.sendEmail(user.Email, subject, body)
}

func (s *NotificationServiceImpl) VerifyOTP(userID uint, code string) (*model.User, error) {
	otpRecord, err := s.repo.Find(userID, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired OTP")
		}
		return nil, err
	}

	if otpRecord.Attempts >= 5 {
		return nil, errors.New("too many attempts")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		// User not found is a system error, don't increment OTP attempts
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		user.IsActive = true
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
	}

	// Increment attempts for this verification attempt
	// If verification succeeds, OTP will be deleted anyway
	if err := s.repo.IncrementAttempts(userID, code); err != nil {
		return nil, err
	}

	// Delete OTP after successful verification
	err = s.repo.Delete(userID, code)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *NotificationServiceImpl) ResendOTP(userID uint) (string, error) {
	canResend, err := s.repo.CanResend(userID)
	if err != nil {
		return "", err
	}
	if !canResend {
		return "", errors.New("please wait before requesting a new OTP")
	}

	code, err := s.GenerateOTP(userID)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}

	if err := s.SendVerificationOTP(user, code); err != nil {
		return "", err
	}
	// UpdateLastSent is not needed here since GenerateOTP already sets LastSent to time.Now()

	return code, nil
}

func (s *NotificationServiceImpl) SendPasswordResetEmail(user *model.User) error {
	code, err := s.GenerateOTP(user.ID)
	if err != nil {
		return err
	}
	subject := "Reset your password"
	body := fmt.Sprintf(`
	<p>Hi %s,</p>
	<p>Please use the following code to reset your password:</p>
	<h3>%s</h3>
	<p>This code expires in 5 minutes.</p>
	`, user.Email, code)

	return s.sendEmail(user.Email, subject, body)
}

// Hàm gửi mail chung
func (s *NotificationServiceImpl) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.sender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUser, s.smtpPass)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}
	log.Printf("Email sent to %s", to)
	return nil
}
