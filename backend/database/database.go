package database

import (
	"todo-list/backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	return DB.AutoMigrate(&models.Todo{}, &models.List{}, &models.Attachment{}, &models.Subtask{})
}
