package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
	JWTSecret string
	JWTIssuer string
	JWTExpiry time.Duration
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found")
	}

	jwtExpiry := 24 * time.Hour
	if value := os.Getenv("JWT_EXPIRY_HOURS"); value != "" {
		if hours, parseErr := strconv.Atoi(value); parseErr == nil && hours > 0 {
			jwtExpiry = time.Duration(hours) * time.Hour
		}
	}

	return Config{
		DBHost:    os.Getenv("DBHOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASS"),
		DBName:    os.Getenv("DB_NAME"),
		DBSSLMode: os.Getenv("DB_SSL_MODE"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTIssuer: getEnvOrDefault("JWT_ISSUER", "donation-system"),
		JWTExpiry: jwtExpiry,
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (cfg Config) ValidateJWT() error {
	if len(cfg.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if cfg.JWTIssuer == "" {
		return fmt.Errorf("JWT_ISSUER must not be empty")
	}
	if cfg.JWTExpiry <= 0 {
		return fmt.Errorf("JWT_EXPIRY_HOURS must be greater than zero")
	}
	return nil
}

func InitDatabase(cfg Config) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	log.Println("Database connected successfully")
	return db
}
