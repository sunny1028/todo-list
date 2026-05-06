package models

import "time"

type FocusSession struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UserID      uint       `json:"user_id" gorm:"index;not null"`
	TodoID      *uint      `json:"todo_id" gorm:"index;default:null"`
	TodoTitle   string     `json:"todo_title" gorm:"-"`
	DurationMin int        `json:"duration_min" gorm:"not null"`
	Completed   bool       `json:"completed" gorm:"default:false"`
	StartedAt   time.Time  `json:"started_at"`
	EndedAt     *time.Time `json:"ended_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
