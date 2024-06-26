package services

import (
	"bbs/dto"
	"bbs/models"
	"bbs/repositories"
)

type ICommentService interface {
	Create(createCommentInput dto.CreateComment, threadId uint, userId uint) (*models.Comment, error)
	FindByThreadId(threadId uint, userId uint) (*[]models.Comment, error)
}

type CommentService struct {
	repository       repositories.ICommentRepository
	threadRepository repositories.IThreadRepository
}

func NewCommentService(repository repositories.ICommentRepository, threadRepository repositories.IThreadRepository) ICommentService {
	return &CommentService{repository: repository, threadRepository: threadRepository}
}

func (s *CommentService) Create(createCommentInput dto.CreateComment, threadId uint, userId uint) (*models.Comment, error) {
	if _, err := s.threadRepository.FindById(threadId); err != nil {
		return nil, err
	}

	newComment := models.Comment{
		Body:     createCommentInput.Body,
		ThreadID: threadId,
		UserID:   userId,
	}

	return s.repository.Create(newComment)
}

func (s *CommentService) FindByThreadId(threadId uint, userId uint) (*[]models.Comment, error) {
	return s.repository.FindByThreadId(threadId, userId)
}
