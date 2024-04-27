package models

import "gorm.io/gorm"

type Thread struct {
	gorm.Model
	Title  string `gorm:"not null"`
	Body   string `gorm:"not null"`
	UserID uint   `gorm:"not null"`
}
