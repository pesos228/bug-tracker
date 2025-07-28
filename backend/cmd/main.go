package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/pesos228/bug-tracker/internal/appmw"
	"github.com/pesos228/bug-tracker/internal/auth"
	"github.com/pesos228/bug-tracker/internal/config"
	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/excel"
	"github.com/pesos228/bug-tracker/internal/handler"
	"github.com/pesos228/bug-tracker/internal/notification"
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

	reportGenerator := excel.NewReportGenerator()
	emailNotifier := notification.NewEmailNotifier(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.From, cfg.AppPublicUrl)

	stateStore := redisstore.NewRedisStateStore(redisClient)
	sessionStore := redisstore.NewRedisSessionStore(redisClient, sessionTTL)
	userStore := psqlstore.NewPsqlUserStore(psqlDb)
	folderStore := psqlstore.NewPsqlFolderStore(psqlDb)
	taskStore := psqlstore.NewPsqlTaskStore(psqlDb)

	authService := service.NewAuthService(&service.AuthServiceDeps{
		AuthClient:   authClient,
		SessionStore: sessionStore,
		StateStore:   stateStore,
		UserStore:    userStore,
		AuthConfig:   &cfg.Auth,
		AppPublicUrl: cfg.AppPublicUrl,
	})
	folderService := service.NewFolderService(folderStore)
	taskService := service.NewTaskService(taskStore, userStore, folderStore, emailNotifier)
	userService := service.NewUserService(userStore, taskStore)
	reportService := service.NewReportService(folderStore, taskStore, reportGenerator)

	authHandler := handler.NewAuthHandler(authService, sessionTTL)
	folderHandler := handler.NewFolderHandler(folderService, reportService)
	taskHandler := handler.NewTaskHandler(taskService)
	userHandler := handler.NewUserHandler(userService)

	authMiddleware := appmw.AuthMiddleware(sessionStore, authClient, authService, userStore)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/auth/login-url", authHandler.HandleLogin)
	r.Get("/auth/callback", authHandler.HandleCallback)
	r.Get("/api/auth/logout-url", authHandler.HandleLogout)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/api/users/me", userHandler.AboutUser)
		r.Get("/api/users/me/stats", userHandler.Stats)
		r.Get("/api/tasks/my", taskHandler.ListUserTasks)
		r.Get("/api/tasks/{id}", taskHandler.Details)
		r.Patch("/api/tasks/{id}/review", taskHandler.UpdateByUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(appmw.AdminOnly)

		r.Get("/api/users", userHandler.Search)

		r.Post("/api/folders", folderHandler.Create)
		r.Get("/api/folders", folderHandler.Search)
		r.Get("/api/folders/{id}", folderHandler.Details)
		r.Get("/api/folders/{id}/tasks", taskHandler.ListByFolder)
		r.Delete("/api/folders/{id}", folderHandler.Delete)

		r.Post("/api/folders/{id}/tasks", taskHandler.Create)
		r.Patch("/api/tasks/{id}", taskHandler.UpdateByAdmin)
		r.Delete("/api/tasks/{id}", taskHandler.Delete)

		r.Get("/api/folders/{id}/reports", folderHandler.Download)
	})

	log.Println("Server started on", cfg.AppPort)
	if err := http.ListenAndServe(cfg.AppPort, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func migrateTables(db *gorm.DB) {
	db.AutoMigrate(domain.User{})
	db.AutoMigrate(domain.Task{})
	db.AutoMigrate(domain.Folder{})
	db.AutoMigrate(domain.Task{})
}
