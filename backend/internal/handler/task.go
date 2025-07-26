package handler

import (
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
	folderID := chi.URLParam(r, "id")
	var newTaskRequest dto.CreateTaskRequest
	if folderID == "" {
		http.Error(w, "Folder id is missing in URL", http.StatusBadRequest)
		return
	}

	if ok := decodeJSON(w, r, &newTaskRequest); !ok {
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
		FolderID:          folderID,
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
	folderID := chi.URLParam(r, "id")
	if folderID == "" {
		http.Error(w, "Folder id is missing in URL", http.StatusBadRequest)
		return
	}
	page := getQueryInt(r.URL.Query(), "page", 1)
	pageSize := getQueryInt(r.URL.Query(), "pageSize", 10)
	checkStatus := getQueryString(r.URL.Query(), "checkStatus", "")
	requestID := getQueryString(r.URL.Query(), "requestID", "")

	tasks, err := t.taskService.SearchByFolderID(r.Context(), &service.SearchTasksByFolderIDParams{
		FolderID:    folderID,
		Page:        page,
		PageSize:    pageSize,
		CheckStatus: checkStatus,
		RequestID:   requestID,
	})

	if err != nil {
		if errors.Is(err, store.ErrFolderNotFound) {
			http.Error(w, fmt.Sprintf("Folder with ID: %s not found", folderID), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSON(w, tasks)
}

func (t *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "Task id is missing in URL", http.StatusBadRequest)
		return
	}

	if err := t.taskService.DeleteByID(r.Context(), taskID); err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			http.Error(w, fmt.Sprintf("Task with ID: %s not found", taskID), http.StatusNotFound)
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
	if ok := decodeJSON(w, r, &taskUpdate); !ok {
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
	if ok := decodeJSON(w, r, &taskUpdate); !ok {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isAdmin := isAdmin(roles)

	if !isAdmin && task.AssigneeID != userID {
		http.Error(w, fmt.Sprintf("user with ID: %s is not assignee task with ID: %s", userID, taskID), http.StatusForbidden)
		return
	}

	switch {
	case strings.EqualFold(view, "full"):
		if !isAdmin {
			http.Error(w, "Forbidden: 'fill' view is not allowed for this user", http.StatusForbidden)
			return
		}

		response := dto.TaskDetailsForAdminResponse{
			ID:                task.ID,
			SoftName:          task.SoftName,
			RequestID:         task.RequestID,
			Description:       task.Description,
			AssigneeID:        task.AssigneeID,
			FolderID:          task.FolderID,
			CheckDate:         task.CheckDate,
			CheckStatus:       task.CheckStatus,
			CheckResult:       task.CheckResult,
			Comment:           task.Comment,
			CreatedAt:         task.CreatedAt,
			TestEnvDateUpdate: task.TestEnvDateUpdate,
		}

		encodeJSON(w, response)
	case strings.EqualFold(view, "short"):
		response := dto.TaskDetailsForUserResponse{
			SoftName:          task.SoftName,
			RequestID:         task.RequestID,
			Description:       task.Description,
			CheckDate:         task.CheckDate,
			CheckStatus:       task.CheckStatus,
			CheckResult:       task.CheckResult,
			Comment:           task.Comment,
			TestEnvDateUpdate: task.TestEnvDateUpdate,
		}

		encodeJSON(w, response)
	default:
		http.Error(w, "unknown view", http.StatusNotFound)
	}
}

func (t *TaskHandler) ListUserTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := appmw.UserIdFromContext(r.Context())
	if !ok || userID == "" {
		http.Error(w, "UserID not found in context", http.StatusInternalServerError)
		return
	}

	page := getQueryInt(r.URL.Query(), "page", 1)
	pageSize := getQueryInt(r.URL.Query(), "pageSize", 10)
	checkStatus := getQueryString(r.URL.Query(), "checkStatus", "")

	tasks, err := t.taskService.SearchByUserID(r.Context(), &service.SearchTasksByUserIDParams{
		AssigneeID:  userID,
		Page:        page,
		PageSize:    pageSize,
		CheckStatus: checkStatus,
	})

	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodeJSON(w, tasks)
}

func isAdmin(s []string) bool {
	for _, role := range s {
		if strings.EqualFold(role, "admin") {
			return true
		}
	}
	return false
}
