package service

import (
	"bbs/internal/dto"
	"bbs/internal/model"
	"bbs/internal/repository"
)

type CommentService struct {
	repository       repository.ICommentRepository
	threadRepository repository.IThreadRepository
}

func NewCommentService(repository repository.ICommentRepository, threadRepository repository.IThreadRepository) ICommentService {
	return &CommentService{repository: repository, threadRepository: threadRepository}
}

func (s *CommentService) Create(createCommentInput dto.CreateComment, threadId uint, userId uint) (*model.Comment, error) {
	if _, err := s.threadRepository.FindById(threadId); err != nil {
		return nil, err
	}

	newComment := model.Comment{
		Body:     createCommentInput.Body,
		ThreadID: threadId,
		UserID:   userId,
	}

	return s.repository.Create(newComment)
}

func (s *CommentService) FindByThreadId(threadId uint, userId uint) (*[]model.Comment, error) {
	return s.repository.FindByThreadId(threadId, userId)
}

func (s *CommentService) FindById(id uint, threadId uint, userId uint) (*model.Comment, error) {
	return s.repository.FindById(id, threadId, userId)
}

func (s *CommentService) Update(updateComment dto.UpdateComment, id uint, threadId uint, userId uint) (*model.Comment, error) {
	targetComment, err := s.repository.FindById(id, threadId, userId)
	if err != nil {
		return nil, err
	}

	targetComment.Body = updateComment.Body

	return s.repository.Update(*targetComment)
}

func (s *CommentService) Delete(id uint, threadId uint, userId uint) error {
	return s.repository.Delete(id, threadId, userId)
}
