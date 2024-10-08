package repository

import (
	"bbs/internal/dto"
	"bbs/internal/model"
	"errors"

	"gorm.io/gorm"
)

type ThreadRepository struct {
	db *gorm.DB
}

func NewThreadRepository(db *gorm.DB) IThreadRepository {
	return &ThreadRepository{db: db}
}

func (r *ThreadRepository) Create(newThread model.Thread) (*model.Thread, error) {
	result := r.db.Create(&newThread)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newThread, nil
}

func (r *ThreadRepository) Update(updateThread model.Thread) (*model.Thread, error) {
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

func (r *ThreadRepository) FindAll(limit int, offset int) (*dto.ThreadListOutput, error) {
	var threadList dto.ThreadListOutput

	// 全レコード数
	r.db.Model(&model.Thread{}).Count(&threadList.Total)

	result := r.db.Limit(limit).Offset(offset * limit).Order("ID desc").Preload("Comments").Find(&threadList.Threads)
	if result.Error != nil {
		return nil, result.Error
	}
	return &threadList, nil
}

func (r *ThreadRepository) FindById(threadId uint) (*model.Thread, error) {
	var thread model.Thread
	result := r.db.First(&thread, "id = ?", threadId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("thread not found")
		}
		return nil, result.Error
	}
	return &thread, nil
}
