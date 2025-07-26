package psqlstore

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/store"
	"gorm.io/gorm"
)

type userStoreImpl struct {
	db *gorm.DB
}

func (u *userStoreImpl) Search(ctx context.Context, params *store.SearchUsersQuery) ([]*domain.User, int64, error) {
	var users []*domain.User
	var count int64

	dbQuery := u.db.WithContext(ctx).Model(&domain.User{})

	if params.FullName != "" {
		words := strings.Fields(params.FullName)

		if len(words) > 0 {
			fullNameExpr := gorm.Expr("CONCAT(first_name, ' ', last_name)")
			for _, word := range words {
				dbQuery = dbQuery.Where("? ILIKE ?", fullNameExpr, "%"+word+"%")
			}
		}
	}

	if err := dbQuery.Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	if count == 0 {
		return []*domain.User{}, 0, nil
	}

	paginatedQuery := dbQuery.Scopes(store.PaginationWithParams(params.Page, params.PageSize)).Find(&users)

	if paginatedQuery.Error != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", paginatedQuery.Error)
	}

	return users, count, nil
}

func (u *userStoreImpl) IsExists(ctx context.Context, userId string) (bool, error) {
	var count int64

	if err := u.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userId).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (u *userStoreImpl) FindAll(ctx context.Context, page, pageSize int, preloads ...store.PreloadOption) ([]*domain.User, int64, error) {
	var users []*domain.User
	var count int64

	query := u.db.WithContext(ctx).Model(&domain.User{})
	query = PreLoad(query, preloads...)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return []*domain.User{}, 0, nil
	}

	paginatedQuery := query.Scopes(store.PaginationWithParams(page, pageSize)).Find(&users)
	if paginatedQuery.Error != nil {
		return nil, 0, paginatedQuery.Error
	}

	return users, count, nil
}

func (u *userStoreImpl) FindById(ctx context.Context, userId string, preloads ...store.PreloadOption) (*domain.User, error) {
	var user domain.User

	query := u.db.WithContext(ctx)
	query = PreLoad(query, preloads...)

	result := query.First(&user, "id = ?", userId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, store.ErrUserNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

func PreLoad(query *gorm.DB, preloads ...store.PreloadOption) *gorm.DB {
	for _, pl := range preloads {
		query = query.Preload(string(pl))
	}
	return query
}

func (u *userStoreImpl) Save(ctx context.Context, user *domain.User) error {
	return u.db.WithContext(ctx).Save(user).Error
}

func NewPsqlUserStore(db *gorm.DB) store.UserStore {
	return &userStoreImpl{db: db}
}
