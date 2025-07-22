package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/google/uuid"
	"github.com/pesos228/bug-tracker/internal/auth"
	"github.com/pesos228/bug-tracker/internal/config"
	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"github.com/pesos228/bug-tracker/internal/store"
	"golang.org/x/oauth2"
)

type AuthService interface {
	PrepareLogin(ctx context.Context) (string, error)
	HandleCallback(ctx context.Context, code, state string) (string, error)
	RefreshToken(ctx context.Context, token string) (*oauth2.Token, error)
	PrepareLogout(ctx context.Context, sessionID string) (string, error)
}

type AuthServiceDeps struct {
	AuthClient   *auth.Client
	SessionStore store.SessionStore
	StateStore   store.StateStore
	UserStore    store.UserStore
	AuthConfig   *config.AuthConfig
	AppPublicUrl string
}

type authServiceImpl struct {
	AuthServiceDeps
}

func (a *authServiceImpl) PrepareLogout(ctx context.Context, sessionID string) (string, error) {
	session, err := a.SessionStore.GetSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, store.ErrSessionNotFound) {
			log.Printf("session %s not found in store, proceeding with Keycloak logout anyway", sessionID)
		} else {
			return "", fmt.Errorf("failed to get session from store: %w", err)
		}
	}

	if err := a.SessionStore.DeleteSession(ctx, sessionID); err != nil {
		return "", fmt.Errorf("failed to delete session from store: %w", err)
	}

	logoutUrl, err := url.Parse(a.AuthClient.LogoutURL(*a.AuthConfig))
	if err != nil {
		return "", fmt.Errorf("could not parse base logout url: %w", err)
	}

	query := logoutUrl.Query()
	query.Set("post_logout_redirect_uri", a.AppPublicUrl+"/")
	if session != nil && session.IdToken != "" {
		query.Set("id_token_hint", session.IdToken)
	}

	logoutUrl.RawQuery = query.Encode()
	return logoutUrl.String(), nil
}

func (a *authServiceImpl) HandleCallback(ctx context.Context, code string, state string) (string, error) {
	storedState, err := a.StateStore.GetState(ctx, state)
	if err != nil || storedState != state {
		if err == store.ErrStateNotFound {
			return "", fmt.Errorf("state not found: %w", err)
		}
		return "", fmt.Errorf("failed to get state: %w", err)
	}

	ouath2Token, err := a.AuthClient.Oauth.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	rawIdToken, ok := ouath2Token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("id_token not found in token response")
	}

	verifiedToken, err := a.AuthClient.OIDC.Verify(ctx, rawIdToken)
	if err != nil {
		return "", fmt.Errorf("failed to verify id_token immediately after exchange: %w", err)
	}

	var claims dto.IdTokenClaims
	if err := verifiedToken.Claims(&claims); err != nil {
		return "", fmt.Errorf("failed to parse token claims: %w", err)
	}

	newUser, err := domain.NewUser(verifiedToken.Subject, claims.Email, claims.GivenName, claims.FamilyName)
	if err != nil {
		return "", err
	}

	if err := a.UserStore.Save(ctx, newUser); err != nil {
		return "", fmt.Errorf("failed to save new user: %w", err)
	}

	sessionId := generateSessionId()

	sessionData := store.SessionData{
		AccessToken:  ouath2Token.AccessToken,
		RefreshToken: ouath2Token.RefreshToken,
		IdToken:      rawIdToken,
	}

	if err := a.SessionStore.SaveSession(ctx, sessionId, &sessionData); err != nil {
		return "", fmt.Errorf("failed to save session: %w", err)
	}

	return sessionId, nil
}

func (a *authServiceImpl) PrepareLogin(ctx context.Context) (string, error) {
	state := generateState()

	if err := a.StateStore.SetState(ctx, state); err != nil {
		return "", fmt.Errorf("failed to set state: %w", err)
	}

	return a.AuthClient.Oauth.AuthCodeURL(state), nil
}

func (a *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	tokenSource := a.AuthClient.Oauth.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
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

func NewAuthService(deps *AuthServiceDeps) AuthService {
	return &authServiceImpl{
		AuthServiceDeps: *deps,
	}
}
