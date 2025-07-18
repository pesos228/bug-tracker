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

type TaskUpdateByAdminRequest struct {
	SoftName          *string    `json:"softName"`
	RequestID         *string    `json:"requestID"`
	Description       *string    `json:"description"`
	TestEnvDateUpdate *time.Time `json:"testEnvDateUpdate"`
	AssigneeID        *string    `json:"assigneeID"`
	FolderID          *string    `json:"folderID"`
	CheckDate         *time.Time `json:"checkDate"`
	CheckStatus       *string    `json:"checkStatus"`
	CheckResult       *string    `json:"checkResult"`
	Comment           *string    `json:"comment"`
}

type TaskUpdateByUserRequest struct {
	CheckStatus *string `json:"checkStatus"`
	CheckResult *string `json:"checkResult"`
	Comment     *string `json:"comment"`
}

type TaskDetailsForAdminResponse struct {
	ID                string     `json:"id"`
	SoftName          string     `json:"softName"`
	RequestID         string     `json:"requestID"`
	Description       string     `json:"description"`
	AssigneeID        string     `json:"assigneeID"`
	FolderID          string     `json:"folderID"`
	TestEnvDateUpdate time.Time  `json:"testEnvDateUpdate"`
	CheckDate         *time.Time `json:"checkDate"`
	CheckStatus       *string    `json:"checkStatus"`
	CheckResult       *string    `json:"checkResult"`
	Comment           *string    `json:"comment"`
	CreatedAt         time.Time  `json:"createdAt"`
}

type TaskDetailsForUserResponse struct {
	SoftName          string     `json:"softName"`
	RequestID         string     `json:"requestID"`
	Description       string     `json:"description"`
	TestEnvDateUpdate time.Time  `json:"testEnvDateUpdate"`
	CheckDate         *time.Time `json:"checkDate"`
	CheckStatus       *string    `json:"checkStatus"`
	CheckResult       *string    `json:"checkResult"`
	Comment           *string    `json:"comment"`
}
