package entity

import "time"

type Product struct {
	ID          uint   `gorm:"primarykey"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	UserID      uint64  `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
