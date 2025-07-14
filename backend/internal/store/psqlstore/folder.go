package psqlstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/store"
	"gorm.io/gorm"
)

type folderStoreImpl struct {
	db *gorm.DB
}

func (f *folderStoreImpl) IsExists(ctx context.Context, folderId string) (bool, error) {
	var count int64

	if err := f.db.WithContext(ctx).Model(domain.Folder{}).Where("id = ?", folderId).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (f *folderStoreImpl) Search(ctx context.Context, page int, pageSize int, query string) ([]*domain.Folder, int64, error) {
	var folders []*domain.Folder
	var count int64

	dbQuery := f.db.WithContext(ctx).Model(&domain.Folder{})

	if query != "" {
		searchQuery := strings.ToLower(query)
		searchPattern := fmt.Sprintf("%%%s%%", searchQuery)
		dbQuery = dbQuery.Where("LOWER(name) LIKE ?", searchPattern)
	}

	if err := dbQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return []*domain.Folder{}, 0, nil
	}

	paginatedQuery := dbQuery.Scopes(store.PaginationWithParams(page, pageSize)).Find(&folders)
	if paginatedQuery.Error != nil {
		return nil, 0, paginatedQuery.Error
	}

	return folders, count, nil
}

func (f *folderStoreImpl) Save(ctx context.Context, folder *domain.Folder) error {
	return f.db.WithContext(ctx).Save(folder).Error
}

func NewPsqlFolderStore(db *gorm.DB) store.FolderStore {
	return &folderStoreImpl{db: db}
}
