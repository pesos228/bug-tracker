package store

import (
	"context"
	"time"

	"github.com/pesos228/bug-tracker/internal/domain"
)

type SessionData struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	IdToken        string `json:"id_token"`
	AbsoluteExpiry int64  `json:"absolute_expiry"`
}

type PreloadOption string

const (
	WithTasks PreloadOption = "Tasks"
)

type StateStore interface {
	SetState(ctx context.Context, state string) error
	GetState(ctx context.Context, state string) (string, error)
	DeleteState(ctx context.Context, state string) error
}

type SessionStore interface {
	SaveSession(ctx context.Context, sessionId string, session *SessionData, ttl ...time.Duration) error
	GetSession(ctx context.Context, sessionId string) (*SessionData, error)
	DeleteSession(ctx context.Context, sessionId string) error
	CheckSession(ctx context.Context, sessionId string) (bool, error)
}

type UserStore interface {
	Save(ctx context.Context, user *domain.User) error
	FindById(ctx context.Context, userId string, preloads ...PreloadOption) (*domain.User, error)
	FindAll(ctx context.Context, page, pageSize int, preloads ...PreloadOption) ([]*domain.User, int64, error)
}

type TaskStore interface {
	Save(ctx context.Context, task *domain.Task) error
	FindById(ctx context.Context, taskId string) (*domain.Task, error)
	FindByUserId(ctx context.Context, page, pageSize int, userId string) ([]*domain.Task, int64, error)
}
