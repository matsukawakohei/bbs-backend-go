package model

import "gorm.io/gorm"

type Thread struct {
	gorm.Model
	Title    string    `gorm:"not null" json:"title"`
	Body     string    `gorm:"not null" json:"body"`
	UserID   uint      `gorm:"not null" json:"userId"`
	Comments []Comment `gorm:"constraint:OnDlete:CASCADE" json:"comments"`
}
