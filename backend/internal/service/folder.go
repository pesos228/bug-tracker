package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"github.com/pesos228/bug-tracker/internal/store"
)

type FolderService interface {
	Save(ctx context.Context, name, userId string) (*dto.FolderCreatedResponse, error)
	Search(ctx context.Context, page, pageSize int, query string) (*dto.FolderSearchResponse, error)
	Delete(ctx context.Context, folderID string) error
	Details(ctx context.Context, folderId string) (*dto.FolderDetailsResponse, error)
}

type folerServiceImpl struct {
	folderStore store.FolderStore
}

func (f *folerServiceImpl) Details(ctx context.Context, folderId string) (*dto.FolderDetailsResponse, error) {
	folder, err := f.folderStore.FindByID(ctx, folderId, store.WithCreator)
	if err != nil {
		if errors.Is(err, store.ErrFolderNotFound) {
			return nil, fmt.Errorf("%w: with ID: %s", err, folderId)
		}
		return nil, fmt.Errorf("db error: %w", err)
	}

	return &dto.FolderDetailsResponse{
		Name:           folder.Name,
		CreatedAt:      folder.CreatedAt,
		AssigneePerson: fmt.Sprintf("%s %s", folder.Creator.LastName, folder.Creator.FirstName),
	}, nil
}

func (f *folerServiceImpl) Delete(ctx context.Context, folderID string) error {
	fodler, err := f.folderStore.FindByID(ctx, folderID)
	if err != nil {
		if errors.Is(err, store.ErrFolderNotFound) {
			return fmt.Errorf("%w: not found with ID: %s", err, folderID)
		}
		return fmt.Errorf("db error: %w", err)
	}

	fodler.Delete()

	if err := f.folderStore.Save(ctx, fodler); err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	return nil
}

func (f *folerServiceImpl) Save(ctx context.Context, name, userId string) (*dto.FolderCreatedResponse, error) {
	newFolder, err := domain.NewFolder(name, userId)
	if err != nil {
		return nil, err
	}
	if err := f.folderStore.Save(ctx, newFolder); err != nil {
		return nil, fmt.Errorf("%s: failed to save folder", err.Error())
	}

	response := &dto.FolderCreatedResponse{
		FolderDataResponse: dto.FolderDataResponse{
			Name:      newFolder.Name,
			CreatedAt: newFolder.CreatedAt,
			Id:        newFolder.ID,
			TaskCount: 0,
		},
	}

	return response, nil
}

func (f *folerServiceImpl) Search(ctx context.Context, page int, pageSize int, query string) (*dto.FolderSearchResponse, error) {
	result, count, err := f.folderStore.Search(ctx, page, pageSize, query)
	if err != nil {
		return nil, fmt.Errorf("%s: error while searching", err.Error())
	}

	data := make([]*dto.FolderDataResponse, len(result))
	for i, folder := range result {
		data[i] = &dto.FolderDataResponse{
			Name:      folder.Name,
			CreatedAt: folder.CreatedAt,
			Id:        folder.ID,
			TaskCount: int(folder.TaskCount),
		}
	}

	return &dto.FolderSearchResponse{
		Data:       data,
		Pagination: store.CalculatePaginationResult(page, pageSize, count),
	}, nil
}

func NewFolderService(folderStore store.FolderStore) FolderService {
	return &folerServiceImpl{folderStore: folderStore}
}
