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

type UpdateTaskParams struct {
	SoftName          *string
	RequestID         *string
	Description       *string
	TestEnvDateUpdate *time.Time
	AssigneeID        *string
	FolderID          *string
	CheckDate         *time.Time
	CheckStatus       *string
	CheckResult       *string
	Comment           *string
}

func (t *Task) validate() error {
	switch {
	case t.SoftName == "":
		return fmt.Errorf("%w: softName is required", ErrValidation)
	case t.RequestID == "":
		return fmt.Errorf("%w: requestId is required", ErrValidation)
	case t.Description == "":
		return fmt.Errorf("%w: description is required", ErrValidation)
	case t.AssigneeID == "":
		return fmt.Errorf("%w: assigneeId is required", ErrValidation)
	case t.CreatorID == "":
		return fmt.Errorf("%w: creatorId is required", ErrValidation)
	case t.FolderID == "":
		return fmt.Errorf("%w: folderId is required", ErrValidation)
	case t.TestEnvDateUpdate.IsZero():
		return fmt.Errorf("%w: testEnvDateUpdate is required", ErrValidation)
	}

	if err := t.CheckStatus.isValid(); err != nil {
		return err
	}
	if err := t.CheckResult.isValid(); err != nil {
		return err
	}

	if t.CheckStatus == NotChecked && t.CheckResult != "" {
		return fmt.Errorf("%w: checkResult must be empty when status is 'not_checked'", ErrValidation)
	}

	if t.CheckStatus != NotChecked && (t.CheckDate == nil || t.CheckDate.IsZero()) {
		return fmt.Errorf("%w: checkDate must be set when status is not 'not_checked'", ErrValidation)
	}

	if t.CheckStatus == Failed && t.CheckResult == Success {
		return fmt.Errorf("%w: checkResult cannot be 'success' when status is 'failed_check'", ErrValidation)
	}

	return nil
}

func NewTask(params *NewTaskParams) (*Task, error) {
	task := &Task{
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
	}

	if err := task.validate(); err != nil {
		return nil, err
	}

	return task, nil
}

func (t *Task) Update(params *UpdateTaskParams) error {
	if params.SoftName != nil {
		t.SoftName = *params.SoftName
	}
	if params.RequestID != nil {
		t.RequestID = *params.RequestID
	}
	if params.Description != nil {
		t.Description = *params.Description
	}
	if params.TestEnvDateUpdate != nil {
		t.TestEnvDateUpdate = *params.TestEnvDateUpdate
	}
	if params.AssigneeID != nil {
		t.AssigneeID = *params.AssigneeID
	}
	if params.FolderID != nil {
		t.FolderID = *params.FolderID
	}
	if params.CheckDate != nil {
		t.CheckDate = params.CheckDate
	}
	if params.CheckStatus != nil {
		t.CheckStatus = CheckStatus(*params.CheckStatus)
	}
	if params.CheckResult != nil {
		t.CheckResult = CheckResult(*params.CheckResult)
	}
	if params.Comment != nil {
		t.Comment = *params.Comment
	}

	return t.validate()
}

func (cs CheckStatus) isValid() error {
	switch cs {
	case NotChecked, Checked, PartiallyChecked, Failed:
		return nil
	default:
		return fmt.Errorf("%w: unknown check status '%s'", ErrValidation, cs)
	}
}

func (cr CheckResult) isValid() error {
	if cr == "" {
		return nil
	}
	switch cr {
	case Success, Failure, Warning:
		return nil
	default:
		return fmt.Errorf("%w: unknown check result '%s'", ErrValidation, cr)
	}
}
