package repositories

import (
	"bbs/internal/models"
	"errors"

	"gorm.io/gorm"
)

type ICommentRepository interface {
	Create(newComment models.Comment) (*models.Comment, error)
	FindByThreadId(threadId uint, userId uint) (*[]models.Comment, error)
	FindById(id uint, threadId uint, userId uint) (*models.Comment, error)
	Update(updateComment models.Comment) (*models.Comment, error)
	Delete(id uint, threadId uint, userId uint) error
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, result.Error
	}

	return &comments, nil
}

func (r *CommentRepository) FindById(id uint, threadId uint, userId uint) (*models.Comment, error) {
	var comment models.Comment
	result := r.db.First(&comment, "id = ? AND thread_id = ? AND user_id = ?", id, threadId, userId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, result.Error
	}

	return &comment, nil
}

func (r *CommentRepository) Update(updateComment models.Comment) (*models.Comment, error) {
	result := r.db.Save(&updateComment)
	if result.Error != nil {
		return nil, result.Error
	}

	return &updateComment, nil
}

func (r *CommentRepository) Delete(id uint, threadId uint, userId uint) error {
	deleteComment, err := r.FindById(id, threadId, userId)
	if err != nil {
		return err
	}
	result := r.db.Delete(&deleteComment)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
