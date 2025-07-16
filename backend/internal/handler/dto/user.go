package dto

type UserPreview struct {
	ID                   string `json:"id"`
	FullName             string `json:"fullName"`
	InProgressTasksCount int    `json:"inProgressTasksCount"`
	CompletedTasksCount  int    `json:"completedTasksCount"`
}

type UserListResponse struct {
	Data       []*UserPreview   `json:"data"`
	Pagination PaginationResult `json:"pagination"`
}
