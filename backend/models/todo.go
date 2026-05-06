package models

import "time"

type Todo struct {
	UserID      uint      `json:"user_id" gorm:"index;not null;default:0"`
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null" binding:"required"`
	Description string    `json:"description"`
	Priority    string    `json:"priority" gorm:"default:'medium'"`
	Effort      string    `json:"effort" gorm:"default:''"`
	Tags        string    `json:"tags" gorm:"default:''"`
	Completed   bool      `json:"completed" gorm:"default:false"`
	Archived    bool      `json:"archived" gorm:"default:false"`
	DueDate     DateOnly  `json:"due_date"`
	Recurrence  string    `json:"recurrence" gorm:"default:''"`
	ListID      uint      `json:"list_id" gorm:"default:0"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SubtaskCount    int `json:"subtask_count" gorm:"-"`
	SubtaskCompleted int `json:"subtask_completed" gorm:"-"`
	FocusMinutes     int `json:"focus_minutes" gorm:"default:0"`
}
