package dto

type CreateComment struct {
	Body string `json:"body" binding:"required"`
}

type UpdateComment struct {
	Body string `json:"body" binding:"required"`
}
