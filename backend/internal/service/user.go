package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"github.com/pesos228/bug-tracker/internal/store"
)

type SearchUsersParams struct {
	Page     int
	PageSize int
	FullName string
}

type UserService interface {
	Search(ctx context.Context, params *SearchUsersParams) (*dto.UserListResponse, error)
}

type userServiceImpl struct {
	userStore store.UserStore
	taskStore store.TaskStore
}

var (
	inProgressCheckStatuses = []domain.CheckStatus{
		domain.NotChecked,
	}

	completedCheckStatuses = []domain.CheckStatus{
		domain.Checked,
		domain.Failed,
		domain.PartiallyChecked,
	}
)

func (u *userServiceImpl) Search(ctx context.Context, params *SearchUsersParams) (*dto.UserListResponse, error) {
	users, count, err := u.userStore.Search(ctx, &store.SearchUsersQuery{
		Page:     params.Page,
		PageSize: params.PageSize,
		FullName: strings.TrimSpace(params.FullName),
	})

	if err != nil {
		return nil, fmt.Errorf("error while searching users: %w", err)
	}

	if len(users) == 0 {
		return &dto.UserListResponse{
			Data:       []*dto.UserPreview{},
			Pagination: store.CalculatePaginationResult(params.Page, params.PageSize, count),
		}, nil
	}

	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	taskCounts, err := u.taskStore.GetTaskCountsForUsers(ctx, userIDs, inProgressCheckStatuses, completedCheckStatuses)
	if err != nil {
		return nil, fmt.Errorf("error while getting task counts: %w", err)
	}

	taskCountsMap := make(map[string]*store.TaskCountResult, len(taskCounts))
	for _, task := range taskCounts {
		taskCountsMap[task.UserID] = task
	}

	data := make([]*dto.UserPreview, len(users))
	for i, user := range users {
		counts := taskCountsMap[user.ID]

		inProgressCount := 0
		completedCount := 0

		if counts != nil {
			inProgressCount = counts.InProgressTasksCount
			completedCount = counts.CompletedTasksCount
		}

		data[i] = &dto.UserPreview{
			ID:                   user.ID,
			FullName:             fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			InProgressTasksCount: inProgressCount,
			CompletedTasksCount:  completedCount,
		}
	}

	pagination := store.CalculatePaginationResult(params.Page, params.PageSize, count)

	return &dto.UserListResponse{
		Data:       data,
		Pagination: pagination,
	}, nil
}

func NewUserService(userStore store.UserStore, taskStore store.TaskStore) UserService {
	return &userServiceImpl{userStore: userStore, taskStore: taskStore}
}
