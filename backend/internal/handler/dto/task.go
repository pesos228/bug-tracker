package dto

import "time"

type CreateTaskRequest struct {
	SoftName          string    `json:"soft_name"`
	RequestId         string    `json:"request_id"`
	Description       string    `json:"description"`
	TestEnvDateUpdate time.Time `json:"test_env_date_update"`
	AssigneeId        string    `json:"assignee_id"`
}
