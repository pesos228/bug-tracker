package store

import (
	"context"
	"time"

	"github.com/pesos228/bug-tracker/internal/domain"
)

type SessionData struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	IdToken        string `json:"id_token"`
	AbsoluteExpiry int64  `json:"absolute_expiry"`
}

type SearchTaskQueryByFolderID struct {
	FolderID    string
	RequestID   string
	Page        int
	PageSize    int
	CheckStatus string
}

type SearchTaskQueryByUserID struct {
	AssigneeID  string
	Page        int
	PageSize    int
	CheckStatus string
}

type SearchUsersQuery struct {
	Page     int
	PageSize int
	FullName string
}

type TaskCountResult struct {
	UserID               string
	InProgressTasksCount int
	CompletedTasksCount  int
}

type FolderSearchResult struct {
	domain.Folder
	TaskCount int64
}

type TasksWithUserInfo struct {
	domain.Task
	FirstName string
	LastName  string
}

type PreloadOption string

const (
	WithTasks PreloadOption = "Tasks"
)

type StateStore interface {
	SetState(ctx context.Context, state string) error
	GetState(ctx context.Context, state string) (string, error)
	DeleteState(ctx context.Context, state string) error
}

type SessionStore interface {
	SaveSession(ctx context.Context, sessionId string, session *SessionData, ttl ...time.Duration) error
	GetSession(ctx context.Context, sessionId string) (*SessionData, error)
	DeleteSession(ctx context.Context, sessionId string) error
	CheckSession(ctx context.Context, sessionId string) (bool, error)
}

type UserStore interface {
	Save(ctx context.Context, user *domain.User) error
	FindById(ctx context.Context, userId string, preloads ...PreloadOption) (*domain.User, error)
	IsExists(ctx context.Context, userId string) (bool, error)
	FindAll(ctx context.Context, page, pageSize int, preloads ...PreloadOption) ([]*domain.User, int64, error)
	Search(ctx context.Context, params *SearchUsersQuery) ([]*domain.User, int64, error)
}

type TaskStore interface {
	Save(ctx context.Context, task *domain.Task) error
	FindById(ctx context.Context, taskId string) (*domain.Task, error)
	FindByUserId(ctx context.Context, page, pageSize int, userId string) ([]*domain.Task, int64, error)
	FindByFolderIdWithUserInfo(ctx context.Context, folderID string) ([]*TasksWithUserInfo, error)
	SearchByFolderID(ctx context.Context, params *SearchTaskQueryByFolderID) ([]*domain.Task, int64, error)
	SearchByUserID(ctx context.Context, params *SearchTaskQueryByUserID) ([]*domain.Task, int64, error)
	DeleteByID(ctx context.Context, taskID string) error
	GetTaskCountsForUsers(ctx context.Context, userIDs []string, inProgressStatuses, completedStatuses []domain.CheckStatus) ([]*TaskCountResult, error)
}

type FolderStore interface {
	Save(ctx context.Context, folder *domain.Folder) error
	Search(ctx context.Context, page, pageSize int, query string) ([]*FolderSearchResult, int64, error)
	IsExists(ctx context.Context, folderId string) (bool, error)
	FindByID(ctx context.Context, folderID string) (*domain.Folder, error)
}
