package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type AuthConfig struct {
	InternalBaselUrl      string
	PublicBaseUrl         string
	ClientId              string
	RedirectUrl           string
	ClientSecret          string
	Realm                 string
	SSOMaxLifespanSeconds int
}

type SmtpConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Config struct {
	Auth         AuthConfig
	Smtp         SmtpConfig
	RedisConfig  redis.Options
	AppPort      string
	AppPublicUrl string
	DatabaseUrl  string
}

func LoadFromEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file, system variables are used")
	}

	RedisDb, err := strconv.Atoi(requireEnv("REDIS_DB"))
	if err != nil {
		log.Fatalf("REDIS_DB must be an integer, got: %v", os.Getenv("REDIS_DB"))
	}

	ssoLifespanSeconds, err := strconv.Atoi(requireEnv("KEYCLOAK_SSO_MAX_LIFESPAN_SECONDS"))
	if err != nil {
		log.Fatalf("KEYCLOAK_SSO_MAX_LIFESPAN_SECONDS must be an integer, got: %v", os.Getenv("KEYCLOAK_SSO_MAX_LIFESPAN_SECONDS"))
	}

	smtpPort, err := strconv.Atoi(requireEnv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("SMTP_PORT must be an integer, got: %v", os.Getenv("SMTP_PORT"))
	}

	return &Config{
		Auth: AuthConfig{
			PublicBaseUrl:         requireEnv("KEYCLOAK_PUBLIC_BASE_URL"),
			InternalBaselUrl:      requireEnv("KEYCLOAK_INTERNAL_BASE_URL"),
			ClientId:              requireEnv("KEYCLOAK_CLIENT_ID"),
			RedirectUrl:           requireEnv("KEYCLOAK_REDIRECT_URL"),
			ClientSecret:          requireEnv("KEYCLOAK_CLIENT_SECRET"),
			Realm:                 requireEnv("KEYCLOAK_REALM"),
			SSOMaxLifespanSeconds: ssoLifespanSeconds,
		},
		RedisConfig: redis.Options{
			Addr:     fmt.Sprintf("%s:%s", requireEnv("REDIS_HOST"), requireEnv("REDIS_PORT")),
			Username: requireEnv("REDIS_USERNAME"),
			Password: requireEnv("REDIS_PASSWORD"),
			DB:       RedisDb,
		},
		Smtp: SmtpConfig{
			Host:     requireEnv("SMTP_HOST"),
			Port:     smtpPort,
			Username: requireEnv("SMTP_USERNAME"),
			Password: requireEnv("SMTP_PASSWORD"),
			From:     requireEnv("SMTP_FROM"),
		},
		AppPort:      requireEnv("APP_PORT"),
		DatabaseUrl:  requireEnv("POSTGRES_URL"),
		AppPublicUrl: requireEnv("APP_PUBLIC_URL"),
	}
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
