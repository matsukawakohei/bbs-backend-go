package dto

import "bbs/internal/model"

type CreateThreadInput struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

type UpdateThreadInput struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

type ThreadListOutput struct {
	Total   int64          `json:"total"`
	Threads []model.Thread `json:"threads"`
}
