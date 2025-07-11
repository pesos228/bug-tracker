package store

import (
	"github.com/pesos228/bug-tracker/internal/handler/dto"
	"gorm.io/gorm"
)

func PaginationWithParams(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func CalculatePaginationResult(page, pageSize int, totalCount int64) dto.PaginationResult {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return dto.PaginationResult{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
