package models

type Subtask struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	TodoID    uint   `json:"todo_id" gorm:"not null;index"`
	Title     string `json:"title" gorm:"not null"`
	Completed bool   `json:"completed" gorm:"default:false"`
}
