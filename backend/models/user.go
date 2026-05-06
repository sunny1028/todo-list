package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UUID         string    `json:"uuid" gorm:"uniqueIndex;not null"`
	Username     string    `json:"username" gorm:"uniqueIndex;default:null"`
	PasswordHash string    `json:"-"`
	HasPassword  bool      `json:"has_password" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
