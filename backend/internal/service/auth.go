package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pesos228/bug-tracker/internal/auth"
	"github.com/pesos228/bug-tracker/internal/store"
	"golang.org/x/oauth2"
)

type AuthService interface {
	PrepareLogin(ctx context.Context) (string, error)
	HandleCallback(ctx context.Context, code, state string) (string, error)
	RefreshToken(ctx context.Context, token string) (*oauth2.Token, error)
}

type authServiceImpl struct {
	authClient   *auth.Client
	sessionStore store.SessionStore
	stateStore   store.StateStore
}

func (a *authServiceImpl) HandleCallback(ctx context.Context, code string, state string) (string, error) {
	storedState, err := a.stateStore.GetState(ctx, state)
	if err != nil || storedState != state {
		if err == store.ErrStateNotFound {
			return "", fmt.Errorf("state not found: %w", err)
		}
		return "", fmt.Errorf("failed to get state: %w", err)
	}

	ouath2Token, err := a.authClient.Oauth.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	rawIdToken, ok := ouath2Token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("id_token not found in token response")
	}

	_, err = a.authClient.OIDC.Verify(ctx, rawIdToken)
	if err != nil {
		return "", fmt.Errorf("failed to verify id_token immediately after exchange: %w", err)
	}

	sessionId := generateSessionId()

	sessionData := store.SessionData{
		AccessToken:  ouath2Token.AccessToken,
		RefreshToken: ouath2Token.RefreshToken,
		IdToken:      rawIdToken,
	}

	if err := a.sessionStore.SaveSession(ctx, sessionId, &sessionData); err != nil {
		return "", fmt.Errorf("failed to save session: %w", err)
	}

	return sessionId, nil
}

func (a *authServiceImpl) PrepareLogin(ctx context.Context) (string, error) {
	state := generateState()

	if err := a.stateStore.SetState(ctx, state); err != nil {
		return "", fmt.Errorf("failed to set state: %w", err)
	}

	return a.authClient.Oauth.AuthCodeURL(state), nil
}

func (a *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	tokenSource := a.authClient.Oauth.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	if newToken.RefreshToken == "" {
		newToken.RefreshToken = refreshToken
	}

	return newToken, nil
}

func generateState() string {
	return uuid.NewString()
}

func generateSessionId() string {
	return uuid.NewString()
}

func NewAuthService(authClient *auth.Client, sessionStore store.SessionStore, stateStore store.StateStore) AuthService {
	return &authServiceImpl{
		authClient:   authClient,
		sessionStore: sessionStore,
		stateStore:   stateStore,
	}
}
