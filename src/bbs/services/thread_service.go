package services

import (
	"bbs/dto"
	"bbs/models"
	"bbs/repositories"
)

type IThreadService interface {
	Create(createThreadInput dto.CreateThreadInput, userId uint) (*models.Thread, error)
	FindAll() (*[]models.Thread, error)
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

func (s *ThreadService) FindAll() (*[]models.Thread, error) {
	return s.repository.FindAll()
}
