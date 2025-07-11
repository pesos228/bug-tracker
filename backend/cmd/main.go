package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pesos228/bug-tracker/internal/appmw"
	"github.com/pesos228/bug-tracker/internal/auth"
	"github.com/pesos228/bug-tracker/internal/config"
	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/handler"
	"github.com/pesos228/bug-tracker/internal/service"
	"github.com/pesos228/bug-tracker/internal/store/psqlstore"
	"github.com/pesos228/bug-tracker/internal/store/redisstore"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.LoadFromEnv()
	ctx := context.Background()

	authClient, err := auth.New(ctx, &cfg.Auth)
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
	}

	redisClient := redis.NewClient(&cfg.RedisConfig)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	psqlDb, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("Connection to DB failed : %v", err)
	}

	migrateTables(psqlDb)

	sessionTTL := time.Duration(cfg.Auth.SSOMaxLifespanSeconds) * time.Second

	stateStore := redisstore.NewRedisStateStore(redisClient)
	sessionStore := redisstore.NewRedisSessionStore(redisClient, sessionTTL)
	userStore := psqlstore.NewPsqlUserStore(psqlDb)

	authService := service.NewAuthService(authClient, sessionStore, stateStore, userStore)

	authHandler := handler.NewAuthHandler(authService, sessionTTL)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/auth/login", authHandler.HandleLogin)
	r.Get("/auth/callback", authHandler.HandleCallback)

	r.Group(func(r chi.Router) {
		r.Use(appmw.AuthMiddleware(sessionStore, authClient, authService))

		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(appmw.KeyUserId)
			userEmail := r.Context().Value(appmw.KeyUserEmail)
			userRoles := r.Context().Value(appmw.KeyUserRoles)

			w.Write([]byte("Protected route accessed!\n"))
			w.Write([]byte("User ID: " + userId.(string) + "\n"))
			w.Write([]byte("User Email: " + userEmail.(string) + "\n"))
			w.Write([]byte("User Roles: " + fmt.Sprintf("%v", userRoles) + "\n"))
		})
	})

	log.Println("Server started on", cfg.AppPort)
	if err := http.ListenAndServe(cfg.AppPort, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func migrateTables(db *gorm.DB) {
	db.AutoMigrate(domain.User{})
	db.AutoMigrate(domain.Task{})
}
