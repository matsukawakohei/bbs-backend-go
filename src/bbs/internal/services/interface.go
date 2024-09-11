package services

import (
	"bbs/internal/dto"
	"bbs/internal/model"
)

type IAuthService interface {
	Signup(name string, email string, password string) error
	Login(email string, password string) (*string, error)
	GetUserFromToken(tokenString string) (*model.User, error)
}

type ICommentService interface {
	Create(createCommentInput dto.CreateComment, threadId uint, userId uint) (*model.Comment, error)
	FindByThreadId(threadId uint, userId uint) (*[]model.Comment, error)
	FindById(id uint, threadId uint, userId uint) (*model.Comment, error)
	Update(updateComment dto.UpdateComment, id uint, threadId uint, userId uint) (*model.Comment, error)
	Delete(id uint, threadId uint, userId uint) error
}

type IThreadService interface {
	Create(createThreadInput dto.CreateThreadInput, userId uint) (*model.Thread, error)
	Update(threadId uint, updateThreadInput dto.UpdateThreadInput, userId uint) (*model.Thread, error)
	Delete(threadId uint, userId uint) error
	FindAll() (*[]model.Thread, error)
	FindById(threadId uint) (*model.Thread, error)
}
