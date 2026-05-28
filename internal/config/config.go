package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBPath         string
	JWTSecret      string
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
	SendGridAPIKey string
	EmailProvider  string // "smtp" or "sendgrid"
	MongoURI       string
	MongoDBName    string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // optional

	port := os.Getenv("PORT")
	if port == "" {
		port = "9779"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "worldcup.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-default-key-please-change" // for dev
	}

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if smtpPort == 0 {
		smtpPort = 587
	}

	emailProvider := os.Getenv("EMAIL_PROVIDER")
	if emailProvider == "" {
		emailProvider = "smtp"
	}

	return &Config{
		Port:           port,
		DBPath:         dbPath,
		JWTSecret:      jwtSecret,
		SMTPHost:       os.Getenv("SMTP_HOST"),
		SMTPPort:       smtpPort,
		SMTPUser:       os.Getenv("SMTP_USER"),
		SMTPPassword:   os.Getenv("SMTP_PASSWORD"),
		SendGridAPIKey: os.Getenv("SENDGRID_API_KEY"),
		EmailProvider:  emailProvider,
		MongoURI:       os.Getenv("MONGO_URI"),
		MongoDBName:    os.Getenv("MONGO_DB_NAME"),
	}, nil
}
