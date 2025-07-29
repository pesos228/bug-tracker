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

	var smtpCfg SmtpConfig
	if getBoolEnv("SMTP_ENABLED") {
		smtpCfg = SmtpConfig{
			Host:     requireEnv("SMTP_HOST"),
			Port:     getIntEnv("SMTP_PORT"),
			Username: requireEnv("SMTP_USERNAME"),
			Password: requireEnv("SMTP_PASSWORD"),
			From:     requireEnv("SMTP_FROM"),
		}
	} else {
		log.Println("SMTP is disabled, email notifications will not be sent")
	}

	return &Config{
		Auth: AuthConfig{
			PublicBaseUrl:         requireEnv("KEYCLOAK_PUBLIC_BASE_URL"),
			InternalBaselUrl:      requireEnv("KEYCLOAK_INTERNAL_BASE_URL"),
			ClientId:              requireEnv("KEYCLOAK_CLIENT_ID"),
			RedirectUrl:           requireEnv("KEYCLOAK_REDIRECT_URL"),
			ClientSecret:          requireEnv("KEYCLOAK_CLIENT_SECRET"),
			Realm:                 requireEnv("KEYCLOAK_REALM"),
			SSOMaxLifespanSeconds: getIntEnv("KEYCLOAK_SSO_MAX_LIFESPAN_SECONDS"),
		},
		RedisConfig: redis.Options{
			Addr:     fmt.Sprintf("%s:%s", requireEnv("REDIS_HOST"), requireEnv("REDIS_PORT")),
			Username: requireEnv("REDIS_USERNAME"),
			Password: requireEnv("REDIS_PASSWORD"),
			DB:       getIntEnv("REDIS_DB"),
		},
		Smtp:         smtpCfg,
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

func getIntEnv(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Fatalf("%s must be an integer, got: %v", key, os.Getenv(key))
	}
	return value
}

func getBoolEnv(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		return false
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("%s must be a boolean, got: %v", key, value)
	}
	return boolValue
}
