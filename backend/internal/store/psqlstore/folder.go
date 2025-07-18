package psqlstore

import (
	"context"
	"fmt"

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

func (f *folderStoreImpl) Search(ctx context.Context, page int, pageSize int, query string) ([]*store.FolderSearchResult, int64, error) {
	var results []*store.FolderSearchResult
	var count int64

	countQuery := f.db.WithContext(ctx).Model(&domain.Folder{})
	if query != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query)
		countQuery = countQuery.Where("name ILIKE ?", searchPattern)
	}

	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return []*store.FolderSearchResult{}, 0, nil
	}

	dataQuery := f.db.WithContext(ctx).Model(&domain.Folder{})
	if query != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query)
		dataQuery = dataQuery.Where("name ILIKE ?", searchPattern)
	}

	err := dataQuery.
		Select("folders.id, folders.name, folders.created_by, folders.created_at, COUNT(tasks.id) as task_count").
		Joins("LEFT JOIN tasks ON tasks.folder_id = folders.id").
		Group("folders.id").
		Order("created_at DESC").
		Scopes(store.PaginationWithParams(page, pageSize)).
		Find(&results).Error

	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func (f *folderStoreImpl) Save(ctx context.Context, folder *domain.Folder) error {
	return f.db.WithContext(ctx).Save(folder).Error
}

func NewPsqlFolderStore(db *gorm.DB) store.FolderStore {
	return &folderStoreImpl{db: db}
}
