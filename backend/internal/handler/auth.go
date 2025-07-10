package handler

import (
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
	http.Redirect(w, r, loginUrl, http.StatusFound)
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

	var redirectUrl string
	if redirectCookie, err := r.Cookie("redirect_after_login"); err == nil {
		redirectUrl = redirectCookie.Value
		http.SetCookie(w, &http.Cookie{Name: "redirect_after_login", Value: "", MaxAge: -1, Path: "/"})
	} else {
		redirectUrl = "/"
	}

	http.Redirect(w, r, redirectUrl, http.StatusFound)
}
