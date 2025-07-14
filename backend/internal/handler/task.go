package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pesos228/bug-tracker/internal/appmw"
	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"github.com/pesos228/bug-tracker/internal/service"
	"github.com/pesos228/bug-tracker/internal/store"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (t *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var newTaskRequest dto.CreateTaskRequest
	if id == "" {
		http.Error(w, "Folder id is missing in URL", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&newTaskRequest); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	creatorId, ok := appmw.UserIdFromContext(r.Context())
	if !ok {
		http.Error(w, "User id not found in context", http.StatusInternalServerError)
		return
	}

	err := t.taskService.Save(r.Context(), &service.CreateFolderParams{
		SoftName:          newTaskRequest.SoftName,
		RequestID:         newTaskRequest.RequestId,
		Description:       newTaskRequest.Description,
		TestEnvDateUpdate: newTaskRequest.TestEnvDateUpdate,
		FolderID:          id,
		AssigneeID:        newTaskRequest.AssigneeId,
		CreatorID:         creatorId,
	})

	if err != nil {
		if errors.Is(err, store.ErrFolderNotFound) || errors.Is(err, store.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
