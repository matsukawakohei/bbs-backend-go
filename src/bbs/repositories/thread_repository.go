package repositories

import (
	"bbs/models"
	"errors"

	"gorm.io/gorm"
)

type IThreadRepository interface {
	Create(newThread models.Thread) (*models.Thread, error)
	FindAll() (*[]models.Thread, error)
	FindById(threadId uint) (*models.Thread, error)
}

type ThreadRepository struct {
	db *gorm.DB
}

func NewThreadRepository(db *gorm.DB) IThreadRepository {
	return &ThreadRepository{db: db}
}

func (r *ThreadRepository) Create(newThread models.Thread) (*models.Thread, error) {
	result := r.db.Create(&newThread)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newThread, nil
}

func (r *ThreadRepository) FindAll() (*[]models.Thread, error) {
	var threads []models.Thread
	result := r.db.Order("ID desc").Find(&threads)
	if result.Error != nil {
		return nil, result.Error
	}
	return &threads, nil
}

func (r *ThreadRepository) FindById(threadId uint) (*models.Thread, error) {
	var thread models.Thread
	result := r.db.First(&thread, "id = ?", threadId)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("thread not found")
		}
		return nil, result.Error
	}
	return &thread, nil
}
