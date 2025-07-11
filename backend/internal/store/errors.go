package store

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrStateNotFound   = errors.New("state not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrTaskNotFound    = errors.New("task not found")
)
