package dto

type CreateThreadInput struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

type UpdateThreadInput struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}
