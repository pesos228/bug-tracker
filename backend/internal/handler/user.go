package handler

import (
	"errors"
	"net/http"

	"github.com/pesos228/bug-tracker/internal/appmw"
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"github.com/pesos228/bug-tracker/internal/service"
	"github.com/pesos228/bug-tracker/internal/store"
)

type UserHandler struct {
	userService service.UserService
}

func (u *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r.URL.Query(), "page", 1)
	pageSize := getQueryInt(r.URL.Query(), "pageSize", 10)
	fullName := getQueryString(r.URL.Query(), "fullName", "")

	users, err := u.userService.Search(r.Context(), &service.SearchUsersParams{
		Page:     page,
		PageSize: pageSize,
		FullName: fullName,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSON(w, users)
}

func (u *UserHandler) AboutUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := appmw.UserIdFromContext(r.Context())
	if !ok || userID == "" {
		http.Error(w, "UserID not found in context", http.StatusInternalServerError)
		return
	}

	firstName, ok := appmw.UserFirstNameFromContext(r.Context())
	if !ok {
		http.Error(w, "FirstName not found in context", http.StatusInternalServerError)
		return
	}

	lastName, ok := appmw.UserLastNameFromContext(r.Context())
	if !ok {
		http.Error(w, "LastName not found in context", http.StatusInternalServerError)
		return
	}

	roles, ok := appmw.UserRolesFromContext(r.Context())
	if !ok || len(roles) == 0 {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	response := dto.UserInfoResponse{
		FirstName: firstName,
		LastName:  lastName,
		IsAdmin:   isAdmin(roles),
	}

	encodeJSON(w, response)
}

func (u *UserHandler) Stats(w http.ResponseWriter, r *http.Request) {
	userID, ok := appmw.UserIdFromContext(r.Context())
	if !ok || userID == "" {
		http.Error(w, "UserID not found in context", http.StatusInternalServerError)
		return
	}

	stats, err := u.userService.GetStats(r.Context(), userID)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSON(w, stats)
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
