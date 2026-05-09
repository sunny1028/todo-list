package models

import "time"

type ListShare struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ListID     uint      `json:"list_id" gorm:"index;not null"`
	Code       string    `json:"code" gorm:"uniqueIndex;not null;size:8"`
	Permission string    `json:"permission" gorm:"not null;default:'view'"`
	CreatedAt  time.Time `json:"created_at"`
}

type ListShareMember struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	ShareID  uint      `json:"share_id" gorm:"index;not null"`
	UserID   uint      `json:"user_id" gorm:"index;not null"`
	JoinedAt time.Time `json:"joined_at"`
}
