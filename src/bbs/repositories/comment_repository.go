package repositories

import (
	"bbs/models"
	"errors"

	"gorm.io/gorm"
)

type ICommentRepository interface {
	Create(newComment models.Comment) (*models.Comment, error)
	FindByThreadId(threadId uint, userId uint) (*[]models.Comment, error)
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

func (r *CommentRepository) FindByThreadId(threadId uint, userId uint) (*[]models.Comment, error) {
	var comments []models.Comment
	result := r.db.Where("thread_id = ? AND user_id = ?", threadId, userId).Find(&comments)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("comment not found")
		}
		return nil, result.Error
	}

	return &comments, nil
}
