package repositories

import (
	"bbs/models"
	"errors"

	"gorm.io/gorm"
)

type IThreadRepository interface {
	Create(newThread models.Thread) (*models.Thread, error)
	Update(updateThread models.Thread) (*models.Thread, error)
	Delete(threadId uint, userId uint) error
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

func (r *ThreadRepository) Update(updateThread models.Thread) (*models.Thread, error) {
	result := r.db.Save(&updateThread)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updateThread, nil
}

func (r *ThreadRepository) Delete(threadId uint, userId uint) error {
	deleteThread, err := r.FindById(threadId)
	if err != nil {
		return err
	}
	result := r.db.Delete(&deleteThread)
	if result.Error != nil {
		return result.Error
	}
	return nil
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("thread not found")
		}
		return nil, result.Error
	}
	return &thread, nil
}
