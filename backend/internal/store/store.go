package store

import (
	"context"
	"time"
)

type SessionData struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	IdToken        string `json:"id_token"`
	AbsoluteExpiry int64  `json:"absolute_expiry"`
}

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
