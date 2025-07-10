package store

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrStateNotFound   = errors.New("state not found")
)
