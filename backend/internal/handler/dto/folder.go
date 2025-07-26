package dto

import "time"

type CreateFolderRequest struct {
	Name string
}

type FolderDataResponse struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Id        string    `json:"id"`
	TaskCount int       `json:"taskCount"`
}

type FolderCreatedResponse struct {
	FolderDataResponse
}

type FolderSearchResponse struct {
	Data       []*FolderDataResponse `json:"data"`
	Pagination PaginationResult      `json:"pagination"`
}

type FolderDetailsResponse struct {
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"createdAt"`
	AssigneePerson string    `json:"assigneePerson"`
}
