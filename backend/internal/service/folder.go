package service

import (
	"context"
	"fmt"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"github.com/pesos228/bug-tracker/internal/store"
)

type FolderService interface {
	Save(ctx context.Context, name, userId string) (*dto.FolderCreatedResponse, error)
	Search(ctx context.Context, page, pageSize int, query string) (*dto.FolderSearchResponse, error)
}

type folerServiceImpl struct {
	folderStore store.FolderStore
}

func (f *folerServiceImpl) Save(ctx context.Context, name, userId string) (*dto.FolderCreatedResponse, error) {
	newFolder, err := domain.NewFolder(name, userId)
	if err != nil {
		return nil, err
	}
	if err := f.folderStore.Save(ctx, newFolder); err != nil {
		return nil, fmt.Errorf("%s: failed to save folder", err.Error())
	}

	return &dto.FolderCreatedResponse{
		FolderDataResponse: mapToFolderDataResponse(newFolder),
	}, nil
}

func (f *folerServiceImpl) Search(ctx context.Context, page int, pageSize int, query string) (*dto.FolderSearchResponse, error) {
	result, count, err := f.folderStore.Search(ctx, page, pageSize, query)
	if err != nil {
		return nil, fmt.Errorf("%s: error while searching", err.Error())
	}

	data := make([]dto.FolderDataResponse, len(result))
	for i, folder := range result {
		data[i] = mapToFolderDataResponse(folder)
	}

	return &dto.FolderSearchResponse{
		Data:       data,
		Pagination: store.CalculatePaginationResult(page, pageSize, count),
	}, nil
}

func NewFolderService(folderStore store.FolderStore) FolderService {
	return &folerServiceImpl{folderStore: folderStore}
}

func mapToFolderDataResponse(folder *domain.Folder) dto.FolderDataResponse {
	return dto.FolderDataResponse{
		Name:      folder.Name,
		Id:        folder.ID,
		CreatedAt: folder.CreatedAt,
	}
}
