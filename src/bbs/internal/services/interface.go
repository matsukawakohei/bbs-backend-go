package services

import (
	"bbs/internal/dto"
	"bbs/internal/models"
)

type IAuthService interface {
	Signup(name string, email string, password string) error
	Login(email string, password string) (*string, error)
	GetUserFromToken(tokenString string) (*models.User, error)
}

type ICommentService interface {
	Create(createCommentInput dto.CreateComment, threadId uint, userId uint) (*models.Comment, error)
	FindByThreadId(threadId uint, userId uint) (*[]models.Comment, error)
	FindById(id uint, threadId uint, userId uint) (*models.Comment, error)
	Update(updateComment dto.UpdateComment, id uint, threadId uint, userId uint) (*models.Comment, error)
	Delete(id uint, threadId uint, userId uint) error
}

type IThreadService interface {
	Create(createThreadInput dto.CreateThreadInput, userId uint) (*models.Thread, error)
	Update(threadId uint, updateThreadInput dto.UpdateThreadInput, userId uint) (*models.Thread, error)
	Delete(threadId uint, userId uint) error
	FindAll() (*[]models.Thread, error)
	FindById(threadId uint) (*models.Thread, error)
}
