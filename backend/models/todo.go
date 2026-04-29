package models

import "time"

type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null" binding:"required"`
	Description string    `json:"description"`
	Priority    string    `json:"priority" gorm:"default:'medium'"`
	Tags        string    `json:"tags" gorm:"default:''"`
	Completed   bool      `json:"completed" gorm:"default:false"`
	DueDate     DateOnly  `json:"due_date"`
	ListID      uint      `json:"list_id" gorm:"default:0"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
