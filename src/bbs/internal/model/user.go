package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string    `gorm:"not null"`
	Email    string    `gorm:"not null;unique"`
	Password string    `gorm:"not null"`
	Threads  []Thread  `gorm:"constrant:OnDelete:CASCADE"`
	Comments []Comment `gorm:"constraint:OnDlete:CASCADE"`
}
