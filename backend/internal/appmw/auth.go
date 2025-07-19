package appmw

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pesos228/bug-tracker/internal/auth"
	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"github.com/pesos228/bug-tracker/internal/service"
	"github.com/pesos228/bug-tracker/internal/store"
)

func AuthMiddleware(sessionStore store.SessionStore,
	authClient *auth.Client,
	authService service.AuthService,
	userStore store.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionId, err := r.Cookie("session_id")
			if err != nil || sessionId == nil {
				redirectToLogin(w, r)
				return
			}

			sessionData, err := sessionStore.GetSession(r.Context(), sessionId.Value)
			if err != nil {
				if err == store.ErrSessionNotFound {
					redirectToLogin(w, r)
					return
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			verifiedToken, err := authClient.OIDC.Verify(r.Context(), sessionData.IdToken)
			if err != nil {
				var tokenExpiredError *oidc.TokenExpiredError
				if errors.As(err, &tokenExpiredError) {
					newVerifiedToken, refreshErr := handleTokenRefresh(r, sessionId.Value, sessionData, authService, sessionStore, authClient)
					if refreshErr != nil {
						redirectToLogin(w, r)
						return
					}
					verifiedToken = newVerifiedToken
					go syncUser(context.Background(), newVerifiedToken, userStore)
				} else {
					redirectToLogin(w, r)
					return
				}
			}

			var claims dto.IdTokenClaims
			if err := verifiedToken.Claims(&claims); err != nil {
				http.Error(w, "Failed to parse token claims", http.StatusInternalServerError)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, KeyUserId, verifiedToken.Subject)
			ctx = context.WithValue(ctx, KeyGivenName, claims.GivenName)
			ctx = context.WithValue(ctx, KeyFamilyName, claims.FamilyName)
			ctx = context.WithValue(ctx, KeyUserRoles, claims.RealmAccess.Roles)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "redirect_after_login",
		Value:    r.URL.Path,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/auth/login", http.StatusFound)
}

func handleTokenRefresh(r *http.Request,
	sessionId string,
	sessionData *store.SessionData,
	authService service.AuthService,
	sessionStore store.SessionStore,
	authClient *auth.Client) (*oidc.IDToken, error) {

	newOAuth2Token, err := authService.RefreshToken(r.Context(), sessionData.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token via service: %w", err)
	}

	newRawIDToken, ok := newOAuth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token not found in refreshed token response")
	}

	sessionData.AccessToken = newOAuth2Token.AccessToken
	sessionData.RefreshToken = newOAuth2Token.RefreshToken
	sessionData.IdToken = newRawIDToken

	expiryTime := time.Unix(sessionData.AbsoluteExpiry, 0)
	remainingTTL := time.Until(expiryTime)

	if remainingTTL <= 0 {
		sessionStore.DeleteSession(r.Context(), sessionId)
		return nil, errors.New("sso session has expired based on absolute expiry time")
	}

	if err := sessionStore.SaveSession(r.Context(), sessionId, sessionData, remainingTTL); err != nil {
		return nil, fmt.Errorf("failed to save refreshed session: %w", err)
	}

	verifiedToken, err := authClient.OIDC.Verify(r.Context(), newRawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify newly refreshed token: %w", err)
	}

	return verifiedToken, nil
}

func syncUser(ctx context.Context, token *oidc.IDToken, userStore store.UserStore) {
	var claims dto.IdTokenClaims
	if err := token.Claims(&claims); err != nil {
		log.Printf("SYNC_USER_ERROR: failed to parse claims: %v", err)
		return
	}

	newUser, err := domain.NewUser(token.Subject, claims.Email, claims.GivenName, claims.FamilyName)
	if err != nil {
		log.Printf("SYNC_USER_ERROR: failed to create new user due to validation: %v", err)
		return
	}
	dbUser, err := userStore.FindById(ctx, token.Subject)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			if err := userStore.Save(ctx, newUser); err != nil {
				log.Printf("SYNC_USER_ERROR: db error: %v", err)
			}
			return
		}
		log.Printf("SYNC_USER_ERROR: db error: %v", err)
		return
	}

	if newUser.ID != dbUser.ID || newUser.Email != dbUser.Email || newUser.FirstName != dbUser.FirstName || newUser.LastName != dbUser.LastName {
		log.Printf("SYNC_USER_INFO: User %s data is outdated. Updating.", dbUser.ID)
		if err := userStore.Save(ctx, newUser); err != nil {
			log.Printf("SYNC_USER_ERROR: db error: %v", err)
		}
	}
}
