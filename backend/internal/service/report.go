package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pesos228/bug-tracker/internal/domain"
	"github.com/pesos228/bug-tracker/internal/store"
)

type TaskReportRow struct {
	SoftName          string
	RequestID         string
	Description       string
	AssigneePerson    string
	TestEnvDateUpdate time.Time
	CheckDate         time.Time
	CheckStatus       domain.CheckStatus
	CheckResult       domain.CheckResult
	Comment           string
}

type ReportGenerator interface {
	Generate(tasks []*TaskReportRow) (*bytes.Buffer, error)
}

type ReportData struct {
	FileName string
	Data     *bytes.Buffer
}

type ReportService interface {
	Create(ctx context.Context, folderID string) (*ReportData, error)
}

type reportServiceImpl struct {
	folderStore     store.FolderStore
	taskStore       store.TaskStore
	reportGenerator ReportGenerator
}

func (r *reportServiceImpl) Create(ctx context.Context, folderID string) (*ReportData, error) {
	folder, err := r.folderStore.FindByID(ctx, folderID)
	if err != nil {
		if errors.Is(err, store.ErrFolderNotFound) {
			return nil, fmt.Errorf("%w: with ID: %s", err, folderID)
		}
		return nil, fmt.Errorf("db error: %w", err)
	}

	tasks, err := r.taskStore.FindByFolderIdWithUserInfo(ctx, folderID)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	taskRows := make([]*TaskReportRow, len(tasks))

	for i, task := range tasks {
		taskRows[i] = &TaskReportRow{
			SoftName:          task.SoftName,
			RequestID:         task.RequestID,
			Description:       task.Description,
			AssigneePerson:    fmt.Sprintf("%s %s", task.LastName, task.FirstName),
			TestEnvDateUpdate: task.TestEnvDateUpdate,
			CheckStatus:       task.CheckStatus,
			CheckResult:       task.CheckResult,
			Comment:           task.Comment,
		}

		if task.CheckDate != nil {
			taskRows[i].CheckDate = *task.CheckDate
		}
	}

	report, err := r.reportGenerator.Generate(taskRows)
	if err != nil {
		return nil, err
	}

	return &ReportData{
		FileName: fmt.Sprintf("%s_%s.xlsx", folder.Name, time.Now().Format("2006-01-02_15-04-05")),
		Data:     report,
	}, nil
}

func NewReportService(folderStore store.FolderStore, taskStore store.TaskStore, generator ReportGenerator) ReportService {
	return &reportServiceImpl{folderStore: folderStore, taskStore: taskStore, reportGenerator: generator}
}
