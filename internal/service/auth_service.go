package service

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"time"

	"github.com/go-kipi/worldcup-2026/internal/config"
	"github.com/go-kipi/worldcup-2026/internal/models"
	"github.com/go-kipi/worldcup-2026/internal/pkg/jwtutil"
	"gorm.io/gorm"
)

type AuthService struct {
	db    *gorm.DB
	cfg   *config.Config
	email *EmailService
}

func NewAuthService(db *gorm.DB, cfg *config.Config, email *EmailService) *AuthService {
	return &AuthService{db: db, cfg: cfg, email: email}
}

func (s *AuthService) SendOTP(email string) error {
	otp := generateOTP(6)

	// Create or reuse an OTP record
	otpRecord := models.OTP{
		Email:     email,
		Code:      otp,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	// Optional: invalidate old OTPs for this email, but for simplicity we'll just insert a new one
	if err := s.db.Create(&otpRecord).Error; err != nil {
		return err
	}

	// For now, in local development, we'll log it directly to the console instead of sending an email.
	log.Printf("================================")
	log.Printf("OTP for %s: %s", email, otp)
	log.Printf("================================")

	if err := s.email.SendOTP(email, otp); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) VerifyOTP(email, code string) (string, error) {
	var otp models.OTP
	err := s.db.Where("email = ? AND code = ? AND expires_at > ?", email, code, time.Now()).Order("created_at desc").First(&otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid or expired OTP")
		}
		return "", err
	}

	// Invalidate the OTP so it can't be used again
	s.db.Delete(&otp)

	// Fetch or create the user
	var user models.User
	err = s.db.Where("email = ?", email).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new user, optionally check if total count < 11 to enforce limits
		user = models.User{Email: email, Name: email, Role: "user"} // basic setup
		if err = s.db.Create(&user).Error; err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	// Return JWT
	return jwtutil.GenerateToken(user.ID, user.Email, user.Role, s.cfg)
}

func generateOTP(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return string(b)
}
