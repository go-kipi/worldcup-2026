package service

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/go-kipi/worldcup-2026/internal/config"
	"github.com/go-kipi/worldcup-2026/internal/models"
	"github.com/go-kipi/worldcup-2026/internal/pkg/jwtutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AuthService struct {
	db    *mongo.Database
	cfg   *config.Config
	email *EmailService
}

func NewAuthService(db *mongo.Database, cfg *config.Config, email *EmailService) *AuthService {
	return &AuthService{db: db, cfg: cfg, email: email}
}

func (s *AuthService) VerifyUserOrDomain(email string) error {
	ctx := context.Background()

	// 1. Check if specific email is allowed
	count, err := s.db.Collection("allowed_users").CountDocuments(ctx, bson.M{"email": email})
	if err == nil && count > 0 {
		return nil
	}

	// 2. Check if domain is allowed
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		domain := parts[1]
		count, err = s.db.Collection("allowed_domains").CountDocuments(ctx, bson.M{
			"domain": bson.M{"$in": []string{domain, "@" + domain}},
		})
		if err == nil && count > 0 {
			return nil
		}
	}

	return errors.New("email or domain not authorized")
}

func (s *AuthService) SendOTP(email string) error {
	if err := s.VerifyUserOrDomain(email); err != nil {
		return err
	}

	otp := generateOTP(6)

	// Create or reuse an OTP record
	otpRecord := models.OTP{
		Email:     email,
		Code:      otp,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
	}

	_, err := s.db.Collection("otps").InsertOne(context.Background(), otpRecord)
	if err != nil {
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
	filter := bson.M{
		"email":      email,
		"code":       code,
		"expires_at": bson.M{"$gt": time.Now()},
	}
	opts := options.FindOne().SetSort(bson.M{"created_at": -1})
	err := s.db.Collection("otps").FindOne(context.Background(), filter, opts).Decode(&otp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("invalid or expired OTP")
		}
		return "", err
	}

	// Invalidate the OTP so it can't be used again
	s.db.Collection("otps").DeleteOne(context.Background(), bson.M{"_id": otp.ID})

	// Fetch or create the user
	var user models.User
	err = s.db.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		// Create new user
		user = models.User{
			Email:     email,
			Name:      email,
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		res, err := s.db.Collection("users").InsertOne(context.Background(), user)
		if err != nil {
			return "", err
		}
		if oid, ok := res.InsertedID.(bson.ObjectID); ok {
			user.ID = oid.Hex()
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
