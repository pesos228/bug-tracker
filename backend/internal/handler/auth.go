package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/pesos228/bug-tracker/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	sessionTTL  time.Duration
}

func NewAuthHandler(authService service.AuthService, sessionTTL time.Duration) *AuthHandler {
	return &AuthHandler{authService: authService, sessionTTL: sessionTTL}
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	loginUrl, err := h.authService.PrepareLogin(r.Context())
	if err != nil {
		log.Printf("Error preparing login: %v", err)
		http.Error(w, "Failed to prepare login", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"login_url": loginUrl})
}

func (h *AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	sessionId, err := h.authService.HandleCallback(r.Context(), code, state)
	if err != nil {
		log.Printf("Error in callback handler: %v", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.sessionTTL.Seconds()),
	})

	http.Redirect(w, r, "/auth-callback-success", http.StatusFound)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionId, err := r.Cookie("session_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}
	}

	logoutUrl, err := h.authService.PrepareLogout(r.Context(), sessionId.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: "", MaxAge: -1, Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"logout_url": logoutUrl})
}
