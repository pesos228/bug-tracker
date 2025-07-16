package dto

import "time"

type CreateTaskRequest struct {
	SoftName          string    `json:"softName"`
	RequestId         string    `json:"requestId"`
	Description       string    `json:"description"`
	TestEnvDateUpdate time.Time `json:"testEnvDateUpdate"`
	AssigneeId        string    `json:"assigneeId"`
}

type TaskPreview struct {
	ID          string    `json:"id"`
	CheckStatus string    `json:"checkStatus"`
	SoftName    string    `json:"softName"`
	RequestID   string    `json:"requestId"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TaskPreviewResponse struct {
	Data       []*TaskPreview   `json:"data"`
	Pagination PaginationResult `json:"pagination"`
}
