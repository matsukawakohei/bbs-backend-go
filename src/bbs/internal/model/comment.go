package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Body     string `gorm:"not null" json:"body"`
	UserID   uint   `gorm:"not null" json:"userId"`
	ThreadID uint   `gorm:"not null" json:"threadId"`
}
