package dto

type CreateComment struct {
	Body     string `json:"body" binding:"required"`
	ThreadId uint   `json:"threadId" binding:"required"`
}
