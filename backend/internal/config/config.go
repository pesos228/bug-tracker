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

type Config struct {
	Auth        AuthConfig
	RedisConfig redis.Options
	AppPort     string
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
		AppPort: requireEnv("APP_PORT"),
	}
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
