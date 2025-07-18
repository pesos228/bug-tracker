package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

	err := t.taskService.Save(r.Context(), &service.CreateTaskParams{
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

func (t *TaskHandler) ListByFolder(w http.ResponseWriter, r *http.Request) {
	folderId := chi.URLParam(r, "id")
	if folderId == "" {
		http.Error(w, "Folder id is missing in URL", http.StatusBadRequest)
		return
	}
	page := getQueryInt(r.URL.Query(), "page", 1)
	pageSize := getQueryInt(r.URL.Query(), "pageSize", 10)
	checkStatus := getQueryString(r.URL.Query(), "checkStatus", "")

	tasks, err := t.taskService.SearchByFolder(r.Context(), &service.SearchTasksParams{
		FolderID:    folderId,
		Page:        page,
		PageSize:    pageSize,
		CheckStatus: checkStatus,
	})

	if err != nil {
		if errors.Is(err, store.ErrFolderNotFound) {
			http.Error(w, fmt.Sprintf("Folder with ID: %s not found", folderId), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, fmt.Sprintf("Failed to Encode DTO: %s", err.Error()), http.StatusInternalServerError)
	}
}

func (t *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "id")
	if taskId == "" {
		http.Error(w, "Task id is missing in URL", http.StatusBadRequest)
		return
	}

	if err := t.taskService.DeleteByID(r.Context(), taskId); err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			http.Error(w, fmt.Sprintf("Task with ID: %s not found", taskId), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TaskHandler) UpdateByAdmin(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "Task id is missing in URL", http.StatusBadRequest)
		return
	}

	var taskUpdate dto.TaskUpdateByAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&taskUpdate); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := t.taskService.UpdateByAdmin(r.Context(), &service.UpdateTaskParams{
		SoftName:          taskUpdate.SoftName,
		RequestID:         taskUpdate.RequestID,
		Description:       taskUpdate.Description,
		TestEnvDateUpdate: taskUpdate.TestEnvDateUpdate,
		AssigneeID:        taskUpdate.AssigneeID,
		FolderID:          taskUpdate.FolderID,
		CheckDate:         taskUpdate.CheckDate,
		CheckStatus:       taskUpdate.CheckStatus,
		CheckResult:       taskUpdate.CheckResult,
		Comment:           taskUpdate.Comment,
		TaskID:            taskID,
	}); err != nil {
		if errors.Is(err, domain.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, store.ErrFolderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, store.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, store.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TaskHandler) UpdateByUser(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "Task id is missing in URL", http.StatusBadRequest)
		return
	}

	var taskUpdate dto.TaskUpdateByUserRequest
	if err := json.NewDecoder(r.Body).Decode(&taskUpdate); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	userID, ok := appmw.UserIdFromContext(r.Context())
	if !ok || userID == "" {
		http.Error(w, "UserID not found in context", http.StatusInternalServerError)
		return
	}

	if err := t.taskService.UpdateByUser(r.Context(), &service.UpdateTaskParams{
		CheckStatus:   taskUpdate.CheckStatus,
		CheckResult:   taskUpdate.CheckResult,
		Comment:       taskUpdate.Comment,
		TaskID:        taskID,
		CurrentUserID: userID,
	}); err != nil {
		if errors.Is(err, domain.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrNotAssignee) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TaskHandler) Details(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "Task id is missing in URL", http.StatusBadRequest)
		return
	}

	view := getQueryString(r.URL.Query(), "view", "")

	roles, ok := appmw.UserRolesFromContext(r.Context())
	if !ok || len(roles) == 0 {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	userID, ok := appmw.UserIdFromContext(r.Context())
	if !ok || userID == "" {
		http.Error(w, "UserID not found in context", http.StatusInternalServerError)
		return
	}

	task, err := t.taskService.GetDetails(r.Context(), taskID, userID)
	if err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrNotAssignee) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isAdmin := isAdmin(roles)

	w.Header().Set("Content-Type", "application/json")
	switch {
	case strings.EqualFold(view, "full"):
		if !isAdmin {
			http.Error(w, "Forbidden: 'fill' view is not allowed for this user", http.StatusForbidden)
			return
		}

		response := dto.TaskDetailsForAdminResponse{
			ID:          task.ID,
			SoftName:    task.SoftName,
			RequestID:   task.RequestID,
			Description: task.Description,
			AssigneeID:  task.AssigneeID,
			FolderID:    task.FolderID,
			CheckDate:   task.CheckDate,
			CheckStatus: task.CheckStatus,
			CheckResult: task.CheckResult,
			Comment:     task.Comment,
			CreatedAt:   task.CreatedAt,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Failed to Encode DTO: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	case strings.EqualFold(view, "short"):
		response := dto.TaskDetailsForUserResponse{
			SoftName:    task.SoftName,
			RequestID:   task.RequestID,
			Description: task.Description,
			CheckDate:   task.CheckDate,
			CheckStatus: task.CheckStatus,
			CheckResult: task.CheckResult,
			Comment:     task.Comment,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Failed to Encode DTO: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "unknown view", http.StatusNotFound)
	}
}

func isAdmin(s []string) bool {
	for _, role := range s {
		if strings.EqualFold(role, "admin") {
			return true
		}
	}
	return false
}
