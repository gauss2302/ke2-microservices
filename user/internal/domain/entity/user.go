package entity

import "time"

type User struct {
	ID        uint   `gorm:"primarykey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
