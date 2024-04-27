package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Body   string `gorm:"not null"`
	UserID uint   `gorm:"not null"`
}
