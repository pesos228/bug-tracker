package redisstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pesos228/bug-tracker/internal/store"
	"github.com/redis/go-redis/v9"
)

type redisStateStore struct {
	client      *redis.Client
	prefixState string
	defaultTTL  time.Duration
}

func (r *redisStateStore) DeleteState(ctx context.Context, state string) error {
	if err := r.client.Del(ctx, r.buildKey(state)).Err(); err != nil {
		return fmt.Errorf("failed to delete state from redis: %w", err)
	}
	return nil
}

func (r *redisStateStore) GetState(ctx context.Context, state string) (string, error) {
	state, err := r.client.Get(ctx, r.buildKey(state)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", store.ErrStateNotFound
		}
		return "", fmt.Errorf("failed to get state from redis: %w", err)
	}
	return state, nil
}

func (r *redisStateStore) SetState(ctx context.Context, state string) error {
	if err := r.client.Set(ctx, r.buildKey(state), state, r.defaultTTL).Err(); err != nil {
		return fmt.Errorf("failed to set state in redis: %w", err)
	}
	return nil
}

func (r *redisStateStore) buildKey(state string) string {
	return fmt.Sprintf("%s:%s", r.prefixState, state)
}

func NewRedisStateStore(client *redis.Client) store.StateStore {
	return &redisStateStore{
		client:      client,
		prefixState: "auth_state",
		defaultTTL:  5 * time.Minute,
	}
}
