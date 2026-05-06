package models

import "time"

type List struct {
	UserID     uint      `json:"user_id" gorm:"index;not null;default:0"`
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Color     string    `json:"color" gorm:"default:'#6366f1'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
