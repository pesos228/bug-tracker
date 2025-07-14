package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pesos228/bug-tracker/internal/appmw"
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"github.com/pesos228/bug-tracker/internal/service"
)

type FolderHandler struct {
	folderService service.FolderService
}

func NewFolderHandler(folderService service.FolderService) *FolderHandler {
	return &FolderHandler{folderService: folderService}
}

func (f *FolderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newFolderRequest dto.CreateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&newFolderRequest); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	userId, ok := appmw.UserIdFromContext(r.Context())
	if !ok {
		http.Error(w, "User id not found in context", http.StatusInternalServerError)
		return
	}

	folder, err := f.folderService.Save(r.Context(), newFolderRequest.Name, userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create new folder: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(folder); err != nil {
		http.Error(w, fmt.Sprintf("Failed to Encode DTO: %s", err.Error()), http.StatusInternalServerError)
	}
}

func (f *FolderHandler) Search(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	query := r.URL.Query().Get("query")
	query = strings.TrimSpace(query)

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	result, err := f.folderService.Search(r.Context(), page, pageSize, query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error while searching: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("Failed to Encode DTO: %s", err.Error()), http.StatusInternalServerError)
	}
}
