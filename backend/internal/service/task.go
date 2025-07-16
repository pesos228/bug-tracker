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

type SearchTasksParams struct {
	FolderID    string
	Page        int
	PageSize    int
	CheckStatus string
}

type TaskService interface {
	Save(ctx context.Context, params *CreateTaskParams) error
	SearchByFolder(ctx context.Context, params *SearchTasksParams) (*dto.TaskPreviewResponse, error)
	DeleteByID(ctx context.Context, taskID string) error
}

type taskServiceImpl struct {
	taskStore   store.TaskStore
	userStore   store.UserStore
	folderStore store.FolderStore
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

func (t *taskServiceImpl) SearchByFolder(ctx context.Context, params *SearchTasksParams) (*dto.TaskPreviewResponse, error) {
	ok, err := t.folderStore.IsExists(ctx, params.FolderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check folder existence: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("%w: with id %s", store.ErrFolderNotFound, params.FolderID)
	}

	tasks, count, err := t.taskStore.Search(ctx, &store.SearchTaskQuery{
		FolderID:    params.FolderID,
		Page:        params.Page,
		PageSize:    params.PageSize,
		CheckStatus: params.CheckStatus,
	})

	if err != nil {
		return nil, fmt.Errorf("error while searching: %w", err)
	}

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

	pagination := store.CalculatePaginationResult(params.Page, params.PageSize, count)

	return &dto.TaskPreviewResponse{
		Data:       data,
		Pagination: pagination,
	}, nil
}

func (t *taskServiceImpl) Save(ctx context.Context, params *CreateTaskParams) error {
	ok, err := t.userStore.IsExists(ctx, params.AssigneeID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !ok {
		return fmt.Errorf("%w: with id %s", store.ErrUserNotFound, params.AssigneeID)
	}

	ok, err = t.folderStore.IsExists(ctx, params.FolderID)
	if err != nil {
		return fmt.Errorf("failed to check folder existence: %w", err)
	}
	if !ok {
		return fmt.Errorf("%w: with id %s", store.ErrFolderNotFound, params.FolderID)
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

func NewTaskService(taskStore store.TaskStore, userStore store.UserStore, folderStore store.FolderStore) TaskService {
	return &taskServiceImpl{taskStore: taskStore, userStore: userStore, folderStore: folderStore}
}
