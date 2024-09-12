package services

import (
	"bbs/internal/dto"
	"bbs/internal/model"
	"bbs/internal/repository"
	"errors"
)

type ThreadService struct {
	repository repository.IThreadRepository
}

func NewThreadService(repository repository.IThreadRepository) IThreadService {
	return &ThreadService{repository: repository}
}

func (s *ThreadService) Create(createThreadInput dto.CreateThreadInput, userId uint) (*model.Thread, error) {
	newThread := model.Thread{
		Title:  createThreadInput.Title,
		Body:   createThreadInput.Body,
		UserID: userId,
	}
	return s.repository.Create(newThread)
}

func (s *ThreadService) Update(threadId uint, updateThreadInput dto.UpdateThreadInput, userId uint) (*model.Thread, error) {
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

func (s *ThreadService) FindAll() (*[]model.Thread, error) {
	return s.repository.FindAll()
}

func (s *ThreadService) FindById(threadId uint) (*model.Thread, error) {
	return s.repository.FindById(threadId)
}
