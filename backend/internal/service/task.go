package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"github.com/pesos228/bug-tracker/internal/store"
)

type CreateTaskParams struct {
	SoftName          string
	RequestID         string
	Description       string
	TestEnvDateUpdate time.Time
	FolderID          string
	AssigneeID        string
	CreatorID         string
}

type SearchTasksByFolderIDParams struct {
	FolderID    string
	Page        int
	PageSize    int
	CheckStatus string
	RequestID   string
}

type SearchTasksByUserIDParams struct {
	AssigneeID  string
	Page        int
	PageSize    int
	CheckStatus string
	RequestID   string
}

type UpdateTaskParams struct {
	SoftName          *string
	RequestID         *string
	Description       *string
	TestEnvDateUpdate *time.Time
	AssigneeID        *string
	FolderID          *string
	CheckDate         *time.Time
	CheckStatus       *string
	CheckResult       *string
	Comment           *string
	TaskID            string
	CurrentUserID     string
}

type TaskDetails struct {
	ID                string
	SoftName          string
	RequestID         string
	Description       string
	AssigneeID        string
	FolderID          string
	TestEnvDateUpdate time.Time
	CheckDate         *time.Time
	CheckStatus       *string
	CheckResult       *string
	Comment           *string
	CreatedAt         time.Time
}

type TaskService interface {
	Save(ctx context.Context, params *CreateTaskParams) error
	SearchByFolderID(ctx context.Context, params *SearchTasksByFolderIDParams) (*dto.TaskPreviewResponse, error)
	SearchByUserID(ctx context.Context, params *SearchTasksByUserIDParams) (*dto.TaskPreviewResponse, error)
	DeleteByID(ctx context.Context, taskID string) error
	UpdateByAdmin(ctx context.Context, params *UpdateTaskParams) error
	UpdateByUser(ctx context.Context, params *UpdateTaskParams) error
	GetDetails(ctx context.Context, taskID, userID string) (*TaskDetails, error)
}

var ErrNotAssignee = errors.New("user is not the assignee")

type taskServiceImpl struct {
	taskStore   store.TaskStore
	userStore   store.UserStore
	folderStore store.FolderStore
}

func (t *taskServiceImpl) SearchByUserID(ctx context.Context, params *SearchTasksByUserIDParams) (*dto.TaskPreviewResponse, error) {
	if err := t.isUserExists(ctx, params.AssigneeID); err != nil {
		return nil, err
	}

	tasks, count, err := t.taskStore.SearchByUserID(ctx, &store.SearchTaskQueryByUserID{
		AssigneeID:  params.AssigneeID,
		Page:        params.Page,
		PageSize:    params.PageSize,
		CheckStatus: params.CheckStatus,
		RequestID:   params.RequestID,
	})

	if err != nil {
		return nil, fmt.Errorf("error while searching: %w", err)
	}

	data := mapTaskToTaskPreview(tasks)

	pagination := store.CalculatePaginationResult(params.Page, params.PageSize, count)

	return &dto.TaskPreviewResponse{
		Data:       data,
		Pagination: pagination,
	}, nil
}

func (t *taskServiceImpl) GetDetails(ctx context.Context, taskID string, userID string) (*TaskDetails, error) {
	task, err := t.taskStore.FindById(ctx, taskID)
	if err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			return nil, fmt.Errorf("%w: not found with ID: %s", err, taskID)
		}
		return nil, fmt.Errorf("db error: %w", err)
	}

	return &TaskDetails{
		ID:                task.ID,
		SoftName:          task.SoftName,
		RequestID:         task.RequestID,
		Description:       task.Description,
		AssigneeID:        task.AssigneeID,
		FolderID:          task.FolderID,
		CheckDate:         task.CheckDate,
		CheckStatus:       (*string)(&task.CheckStatus),
		CheckResult:       (*string)(&task.CheckResult),
		Comment:           &task.Comment,
		CreatedAt:         task.CreatedAt,
		TestEnvDateUpdate: task.TestEnvDateUpdate,
	}, nil
}

func (t *taskServiceImpl) UpdateByUser(ctx context.Context, params *UpdateTaskParams) error {
	task, err := t.taskStore.FindById(ctx, params.TaskID)
	if err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			return fmt.Errorf("%w: with ID %s", err, params.TaskID)
		}
		return fmt.Errorf("db error: %w", err)
	}

	if task.AssigneeID != params.CurrentUserID {
		return fmt.Errorf("%w: task with ID: %s", ErrNotAssignee, params.TaskID)
	}

	now := time.Now().UTC()
	domainParams := &domain.UpdateTaskParams{
		CheckStatus: params.CheckStatus,
		CheckResult: params.CheckResult,
		Comment:     params.Comment,
		CheckDate:   &now,
	}

	if err := task.Update(domainParams); err != nil {
		return err
	}

	if err := t.taskStore.Save(ctx, task); err != nil {
		return fmt.Errorf("error while updating task: %w", err)
	}

	return nil
}

func (t *taskServiceImpl) UpdateByAdmin(ctx context.Context, params *UpdateTaskParams) error {
	task, err := t.taskStore.FindById(ctx, params.TaskID)
	if err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			return fmt.Errorf("%w: with ID %s", err, params.TaskID)
		}
		return fmt.Errorf("db error: %w", err)
	}

	if params.AssigneeID != nil {
		if err := t.isUserExists(ctx, *params.AssigneeID); err != nil {
			return err
		}
	}

	if params.FolderID != nil {
		if err := t.isFolderExists(ctx, *params.FolderID); err != nil {
			return err
		}
	}

	domainParams := &domain.UpdateTaskParams{
		SoftName:          params.SoftName,
		RequestID:         params.RequestID,
		Description:       params.Description,
		TestEnvDateUpdate: params.TestEnvDateUpdate,
		AssigneeID:        params.AssigneeID,
		FolderID:          params.FolderID,
		CheckDate:         params.CheckDate,
		CheckStatus:       params.CheckStatus,
		CheckResult:       params.CheckResult,
		Comment:           params.Comment,
	}

	if err := task.Update(domainParams); err != nil {
		return err
	}

	if err := t.taskStore.Save(ctx, task); err != nil {
		return fmt.Errorf("error while updating task: %w", err)
	}

	return nil
}

func (t *taskServiceImpl) DeleteByID(ctx context.Context, taskID string) error {
	if err := t.taskStore.DeleteByID(ctx, taskID); err != nil {
		if !errors.Is(err, store.ErrTaskNotFound) {
			return fmt.Errorf("db error: %w", err)
		}
		return err
	}
	return nil
}

func (t *taskServiceImpl) SearchByFolderID(ctx context.Context, params *SearchTasksByFolderIDParams) (*dto.TaskPreviewResponse, error) {
	if err := t.isFolderExists(ctx, params.FolderID); err != nil {
		return nil, err
	}

	tasks, count, err := t.taskStore.SearchByFolderID(ctx, &store.SearchTaskQueryByFolderID{
		FolderID:    params.FolderID,
		RequestID:   params.RequestID,
		Page:        params.Page,
		PageSize:    params.PageSize,
		CheckStatus: params.CheckStatus,
	})

	if err != nil {
		return nil, fmt.Errorf("error while searching: %w", err)
	}

	data := mapTaskToTaskPreview(tasks)

	pagination := store.CalculatePaginationResult(params.Page, params.PageSize, count)

	return &dto.TaskPreviewResponse{
		Data:       data,
		Pagination: pagination,
	}, nil
}

func (t *taskServiceImpl) Save(ctx context.Context, params *CreateTaskParams) error {
	if err := t.isUserExists(ctx, params.AssigneeID); err != nil {
		return err
	}
	if err := t.isFolderExists(ctx, params.FolderID); err != nil {
		return err
	}

	newTask, err := domain.NewTask(&domain.NewTaskParams{
		SoftName:          params.SoftName,
		RequestID:         params.RequestID,
		Description:       params.Description,
		AssigneeID:        params.AssigneeID,
		CreatorID:         params.CreatorID,
		FolderID:          params.FolderID,
		TestEnvDateUpdate: params.TestEnvDateUpdate,
	})

	if err != nil {
		return err
	}

	if err := t.taskStore.Save(ctx, newTask); err != nil {
		return fmt.Errorf("db error while saving: %s", err.Error())
	}

	return nil
}

func (t *taskServiceImpl) isUserExists(ctx context.Context, userID string) error {
	ok, err := t.userStore.IsExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !ok {
		return fmt.Errorf("%w: with id %s", store.ErrUserNotFound, userID)
	}
	return nil
}

func (t *taskServiceImpl) isFolderExists(ctx context.Context, folderID string) error {
	ok, err := t.folderStore.IsExists(ctx, folderID)
	if err != nil {
		return fmt.Errorf("failed to check folder existence: %w", err)
	}
	if !ok {
		return fmt.Errorf("%w: with id %s", store.ErrFolderNotFound, folderID)
	}
	return nil
}

func mapTaskToTaskPreview(tasks []*domain.Task) []*dto.TaskPreview {
	data := make([]*dto.TaskPreview, len(tasks))
	for i, task := range tasks {
		data[i] = &dto.TaskPreview{
			ID:          task.ID,
			CheckStatus: string(task.CheckStatus),
			SoftName:    task.SoftName,
			RequestID:   task.RequestID,
			Description: task.Description,
			CreatedAt:   task.CreatedAt,
		}
	}
	return data
}

func NewTaskService(taskStore store.TaskStore, userStore store.UserStore, folderStore store.FolderStore) TaskService {
	return &taskServiceImpl{taskStore: taskStore, userStore: userStore, folderStore: folderStore}
}
