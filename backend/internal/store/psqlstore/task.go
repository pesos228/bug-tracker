package psqlstore

import (
	"context"
	"errors"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/store"
	"gorm.io/gorm"
)

type taskStoreImpl struct {
	db *gorm.DB
}

func (t *taskStoreImpl) FindById(ctx context.Context, taskId string) (*domain.Task, error) {
	var task domain.Task
	result := t.db.WithContext(ctx).First(&task, "id = ?", taskId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, store.ErrTaskNotFound
		}
		return nil, result.Error
	}

	return &task, nil
}

func (t *taskStoreImpl) FindByUserId(ctx context.Context, page, pageSize int, userId string) ([]*domain.Task, int64, error) {
	var tasks []*domain.Task
	var count int64

	query := t.db.WithContext(ctx).Model(&domain.Task{}).Where("user_id = ?", userId)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return nil, 0, nil
	}

	if err := query.Scopes(store.PaginationWithParams(page, pageSize)).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}

func (t *taskStoreImpl) Save(ctx context.Context, task *domain.Task) error {
	panic("unimplemented")
}

func NewPsqlTaskStore(db *gorm.DB) store.TaskStore {
	return &taskStoreImpl{db: db}
}
