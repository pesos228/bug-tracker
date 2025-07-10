package redisstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pesos228/bug-tracker/internal/store"
	"github.com/redis/go-redis/v9"
)

type redisSessionStore struct {
	client      *redis.Client
	prefixState string
	defaultTTL  time.Duration
}

func (r *redisSessionStore) CheckSession(ctx context.Context, sessionId string) (bool, error) {
	result, err := r.client.Exists(ctx, r.buildKey(sessionId)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check session in redis: %w", err)
	}

	return result == 1, nil
}

func (r *redisSessionStore) DeleteSession(ctx context.Context, sessionId string) error {
	if err := r.client.Del(ctx, r.buildKey(sessionId)).Err(); err != nil {
		return fmt.Errorf("failed to delete session from redis: %w", err)
	}
	return nil
}

func (r *redisSessionStore) GetSession(ctx context.Context, sessionId string) (*store.SessionData, error) {
	data, err := r.client.Get(ctx, r.buildKey(sessionId)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, store.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session from redis: %w", err)
	}
	var session store.SessionData
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}
	return &session, nil
}

func (r *redisSessionStore) SaveSession(ctx context.Context, sessionId string, session *store.SessionData, ttl ...time.Duration) error {
	var currentTTL time.Duration
	if len(ttl) > 0 {
		currentTTL = ttl[0]
	} else {
		session.AbsoluteExpiry = time.Now().Add(r.defaultTTL).Unix()
		currentTTL = r.defaultTTL
	}

	jsonData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}
	if err := r.client.Set(ctx, r.buildKey(sessionId), string(jsonData), currentTTL).Err(); err != nil {
		return fmt.Errorf("failed to save session to redis: %w", err)
	}
	return nil
}

func (r *redisSessionStore) buildKey(sessionId string) string {
	return fmt.Sprintf("%s:%s", r.prefixState, sessionId)
}

func NewRedisSessionStore(rds *redis.Client, defaultTTL time.Duration) store.SessionStore {
	return &redisSessionStore{
		client:      rds,
		prefixState: "session",
		defaultTTL:  defaultTTL,
	}
}
