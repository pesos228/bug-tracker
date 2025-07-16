package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CheckStatus string
type CheckResult string

const (
	NotChecked       CheckStatus = "not_checked"
	Checked          CheckStatus = "checked"
	PartiallyChecked CheckStatus = "partially_checked"
	Failed           CheckStatus = "failed_check"
)

const (
	Success CheckResult = "success"
	Failure CheckResult = "failure"
	Warning CheckResult = "warning"
)

type Task struct {
	BaseModel
	SoftName          string      `gorm:"type:varchar(255)"`
	RequestID         string      `gorm:"type:varchar(255)"`
	Description       string      `gorm:"type:text"`
	AssigneeID        string      `gorm:"type:uuid;not null"`
	CreatorID         string      `gorm:"type:uuid;not null"`
	FolderID          string      `gorm:"type:uuid;not null"`
	TestEnvDateUpdate time.Time   `gorm:"type:date"`
	CheckDate         *time.Time  `gorm:"type:date"`
	CheckStatus       CheckStatus `gorm:"type:varchar(20);not null;default:'not_checked'"`
	CheckResult       CheckResult `gorm:"type:varchar(20)"`
	Comment           string      `gorm:"type:text"`
	CreatedAt         time.Time   `gorm:"type:timestamptz;not null"`
}

type NewTaskParams struct {
	SoftName          string
	RequestID         string
	Description       string
	AssigneeID        string
	CreatorID         string
	FolderID          string
	TestEnvDateUpdate time.Time
}

func (n *NewTaskParams) validate() error {
	switch {
	case n.SoftName == "":
		return fmt.Errorf("%w: softName is required", ErrValidation)
	case n.RequestID == "":
		return fmt.Errorf("%w: requestId is required", ErrValidation)
	case n.Description == "":
		return fmt.Errorf("%w: description is required", ErrValidation)
	case n.AssigneeID == "":
		return fmt.Errorf("%w: assigneeId is required", ErrValidation)
	case n.CreatorID == "":
		return fmt.Errorf("%w: creatorId is required", ErrValidation)
	case n.FolderID == "":
		return fmt.Errorf("%w: folderId is required", ErrValidation)
	case n.TestEnvDateUpdate.IsZero():
		return fmt.Errorf("%w: testEnvDateUpdate is required", ErrValidation)
	default:
		return nil
	}
}

func NewTask(params *NewTaskParams) (*Task, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	return &Task{
		BaseModel: BaseModel{
			uuid.NewString(),
		},
		SoftName:          params.SoftName,
		RequestID:         params.RequestID,
		Description:       params.Description,
		AssigneeID:        params.AssigneeID,
		CreatorID:         params.CreatorID,
		FolderID:          params.FolderID,
		CheckStatus:       NotChecked,
		TestEnvDateUpdate: params.TestEnvDateUpdate,
		CreatedAt:         time.Now().UTC(),
	}, nil
}
