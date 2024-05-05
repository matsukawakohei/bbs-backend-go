package repositories

import (
	"bbs/models"

	"gorm.io/gorm"
)

type IThreadRepository interface {
	Create(newThread models.Thread) (*models.Thread, error)
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
