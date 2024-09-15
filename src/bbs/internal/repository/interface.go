package repository

import (
	"bbs/internal/dto"
	"bbs/internal/model"
)

type IAuthRepository interface {
	CreateUser(user model.User) error
	FindUser(email string) (*model.User, error)
}

type IThreadRepository interface {
	Create(newThread model.Thread) (*model.Thread, error)
	Update(updateThread model.Thread) (*model.Thread, error)
	Delete(threadId uint, userId uint) error
	FindAll(limit int, offset int) (*dto.ThreadListOutput, error)
	FindById(threadId uint) (*model.Thread, error)
}

type ICommentRepository interface {
	Create(newComment model.Comment) (*model.Comment, error)
	FindByThreadId(threadId uint, userId uint) (*[]model.Comment, error)
	FindById(id uint, threadId uint, userId uint) (*model.Comment, error)
	Update(updateComment model.Comment) (*model.Comment, error)
	Delete(id uint, threadId uint, userId uint) error
}
