package models

type Subtask struct {
	UserID     uint   `json:"user_id" gorm:"index;not null;default:0"`
	ID        uint   `json:"id" gorm:"primaryKey"`
	TodoID    uint   `json:"todo_id" gorm:"not null;index"`
	Title     string `json:"title" gorm:"not null"`
	Completed bool   `json:"completed" gorm:"default:false"`
	SortOrder int    `json:"sort_order" gorm:"default:0"`
}
