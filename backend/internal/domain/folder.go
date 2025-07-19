package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Folder struct {
	BaseModel
	Name      string     `gorm:"type:VARCHAR(255);not null"`
	CreatedBy string     `gorm:"type:uuid;not null"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null"`
	DeletedAt *time.Time `gorm:"type:timestamptz"`
}

func NewFolder(name, userId string) (*Folder, error) {
	if name == "" || userId == "" {
		return nil, fmt.Errorf("%w: name or user id is empty", ErrValidation)
	}

	return &Folder{
		BaseModel: BaseModel{
			ID: uuid.NewString(),
		},
		Name:      name,
		CreatedBy: userId,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (f *Folder) Delete() {
	now := time.Now().UTC()
	f.DeletedAt = &now
}
