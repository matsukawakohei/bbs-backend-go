package services

import (
	"bbs/internal/dto"
	"bbs/internal/models"
	"bbs/internal/repositories"
	"errors"
)

type IThreadService interface {
	Create(createThreadInput dto.CreateThreadInput, userId uint) (*models.Thread, error)
	Update(threadId uint, updateThreadInput dto.UpdateThreadInput, userId uint) (*models.Thread, error)
	Delete(threadId uint, userId uint) error
	FindAll() (*[]models.Thread, error)
	FindById(threadId uint) (*models.Thread, error)
}

type ThreadService struct {
	repository repositories.IThreadRepository
}

func NewThreadService(repository repositories.IThreadRepository) IThreadService {
	return &ThreadService{repository: repository}
}

func (s *ThreadService) Create(createThreadInput dto.CreateThreadInput, userId uint) (*models.Thread, error) {
	newThread := models.Thread{
		Title:  createThreadInput.Title,
		Body:   createThreadInput.Body,
		UserID: userId,
	}
	return s.repository.Create(newThread)
}

func (s *ThreadService) Update(threadId uint, updateThreadInput dto.UpdateThreadInput, userId uint) (*models.Thread, error) {
	targetThread, err := s.FindById(threadId)
	if err != nil {
		return nil, err
	}

	if targetThread.UserID != userId {
		return nil, errors.New("user is not thread owner")
	}

	if updateThreadInput.Title != nil {
		targetThread.Title = *updateThreadInput.Title
	}

	if updateThreadInput.Body != nil {
		targetThread.Body = *updateThreadInput.Body
	}

	return s.repository.Update(*targetThread)
}

func (s *ThreadService) Delete(threadId uint, userId uint) error {
	targetThread, err := s.FindById(threadId)
	if err != nil {
		return err
	}

	if targetThread.UserID != userId {
		return errors.New("user is not thread owner")
	}

	return s.repository.Delete(threadId, userId)
}

func (s *ThreadService) FindAll() (*[]models.Thread, error) {
	return s.repository.FindAll()
}

func (s *ThreadService) FindById(threadId uint) (*models.Thread, error) {
	return s.repository.FindById(threadId)
}
