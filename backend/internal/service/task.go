package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/store"
)

type CreateFolderParams struct {
	SoftName          string
	RequestID         string
	Description       string
	TestEnvDateUpdate time.Time
	FolderID          string
	AssigneeID        string
	CreatorID         string
}

type TaskService interface {
	Save(ctx context.Context, params *CreateFolderParams) error
}

type taskServiceImpl struct {
	taskStore   store.TaskStore
	userStore   store.UserStore
	folderStore store.FolderStore
}

func (t *taskServiceImpl) Save(ctx context.Context, params *CreateFolderParams) error {
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
