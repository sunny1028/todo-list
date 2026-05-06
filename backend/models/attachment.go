package models

import "time"

type Attachment struct {
	UserID     uint      `json:"user_id" gorm:"index;not null;default:0"`
	ID        uint      `json:"id" gorm:"primaryKey"`
	TodoID    uint      `json:"todo_id" gorm:"not null;index"`
	Filename  string    `json:"filename" gorm:"not null"`
	Filepath  string    `json:"-" gorm:"not null"`
	MimeType  string    `json:"mime_type"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
