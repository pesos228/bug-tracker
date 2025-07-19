package psqlstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/store"
	"gorm.io/gorm"
)

type taskStoreImpl struct {
	db *gorm.DB
}

func (t *taskStoreImpl) SearchByUserID(ctx context.Context, params *store.SearchTaskQueryByUserID) ([]*domain.Task, int64, error) {
	var tasks []*domain.Task
	var count int64

	dbQuery := t.db.WithContext(ctx).Model(&domain.Task{}).Where("assignee_id = ?", params.AssigneeID).
		Joins("JOIN folders f ON tasks.folder_id = f.id").
		Where("f.deleted_at IS NULL")

	if params.CheckStatus != "" {
		dbQuery = dbQuery.Where("check_status = ?", params.CheckStatus)
	}

	if err := dbQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return []*domain.Task{}, 0, nil
	}

	paginatedQuery := dbQuery.Order("created_at DESC").Scopes(store.PaginationWithParams(params.Page, params.PageSize)).Find(&tasks)
	if paginatedQuery.Error != nil {
		return nil, 0, paginatedQuery.Error
	}

	return tasks, count, nil
}

func (t *taskStoreImpl) GetTaskCountsForUsers(ctx context.Context, userIDs []string, inProgressStatuses, completedStatuses []domain.CheckStatus) ([]*store.TaskCountResult, error) {
	var tasksCount []*store.TaskCountResult

	inProgressExpr := gorm.Expr("COUNT(CASE WHEN check_status IN (?) THEN 1 END)", inProgressStatuses)
	completedExpr := gorm.Expr("COUNT(CASE WHEN check_status IN (?) THEN 1 END)", completedStatuses)

	err := t.db.WithContext(ctx).Model(&domain.Task{}).
		Select("assignee_id as user_id, ? as in_progress_tasks_count, ? as completed_tasks_count", inProgressExpr, completedExpr).
		Where("assignee_id IN (?)", userIDs).
		Group("assignee_id").Find(&tasksCount).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get task counts for users: %w", err)
	}

	return tasksCount, nil
}

func (t *taskStoreImpl) DeleteByID(ctx context.Context, taskID string) error {
	result := t.db.WithContext(ctx).Delete(&domain.Task{}, "id = ?", taskID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return store.ErrTaskNotFound
	}

	return nil
}

func (t *taskStoreImpl) SearchByFolderID(ctx context.Context, params *store.SearchTaskQueryByFolderID) ([]*domain.Task, int64, error) {
	var tasks []*domain.Task
	var count int64

	dbQuery := t.db.WithContext(ctx).Model(&domain.Task{}).Where("folder_id = ?", params.FolderID).
		Joins("JOIN folders f ON tasks.folder_id = f.id").
		Where("f.deleted_at IS NULL")

	if params.CheckStatus != "" {
		dbQuery = dbQuery.Where("check_status = ?", params.CheckStatus)
	}

	if params.RequestID != "" {
		searchPattern := fmt.Sprintf("%%%s%%", params.RequestID)
		dbQuery = dbQuery.Where("request_id ILIKE ?", searchPattern)
	}

	if err := dbQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return []*domain.Task{}, 0, nil
	}

	paginatedQuery := dbQuery.Order("created_at DESC").Scopes(store.PaginationWithParams(params.Page, params.PageSize)).Find(&tasks)
	if paginatedQuery.Error != nil {
		return nil, 0, paginatedQuery.Error
	}

	return tasks, count, nil
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
	return t.db.WithContext(ctx).Save(task).Error
}

func NewPsqlTaskStore(db *gorm.DB) store.TaskStore {
	return &taskStoreImpl{db: db}
}
