package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pesos228/bug-tracker/internal/appmw"
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"github.com/pesos228/bug-tracker/internal/service"
	"github.com/pesos228/bug-tracker/internal/store"
)

type FolderHandler struct {
	folderService service.FolderService
	reportService service.ReportService
}

func NewFolderHandler(folderService service.FolderService, reportService service.ReportService) *FolderHandler {
	return &FolderHandler{folderService: folderService, reportService: reportService}
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

	w.WriteHeader(http.StatusCreated)
	encodeJSON(w, folder)
}

func (f *FolderHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := getQueryString(r.URL.Query(), "query", "")
	query = strings.TrimSpace(query)

	page := getQueryInt(r.URL.Query(), "page", 1)
	pageSize := getQueryInt(r.URL.Query(), "pageSize", 10)

	result, err := f.folderService.Search(r.Context(), page, pageSize, query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error while searching: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	encodeJSON(w, result)
}

func (f *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	folderID := chi.URLParam(r, "id")
	if folderID == "" {
		http.Error(w, "Folder id is missing in URL", http.StatusBadRequest)
		return
	}

	if err := f.folderService.Delete(r.Context(), folderID); err != nil {
		if errors.Is(err, store.ErrFolderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (f *FolderHandler) Download(w http.ResponseWriter, r *http.Request) {
	folderID := chi.URLParam(r, "id")
	if folderID == "" {
		http.Error(w, "Folder id is missing in URL", http.StatusBadRequest)
		return
	}

	report, err := f.reportService.Create(r.Context(), folderID)
	if err != nil {
		if errors.Is(err, store.ErrFolderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", report.FileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", report.Data.Len()))

	w.Write(report.Data.Bytes())
}
