package repositories

import (
	"bbs/models"

	"gorm.io/gorm"
)

type ICommentRepository interface {
	Create(newComment models.Comment) (*models.Comment, error)
}

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) ICommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(newComment models.Comment) (*models.Comment, error) {
	result := r.db.Create(&newComment)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newComment, nil
}
