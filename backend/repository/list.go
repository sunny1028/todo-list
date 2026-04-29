package repository

import (
	"todo-list/backend/database"
	"todo-list/backend/models"
)

func FindAllLists() ([]models.List, error) {
	var lists []models.List
	err := database.DB.Order("created_at ASC").Find(&lists).Error
	return lists, err
}

func FindListByID(id uint) (*models.List, error) {
	var list models.List
	err := database.DB.First(&list, id).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func CreateList(list *models.List) error {
	return database.DB.Create(list).Error
}

func UpdateList(list *models.List) error {
	return database.DB.Save(list).Error
}

func DeleteList(id uint) error {
	// Move todos in this list to no-list (0)
	database.DB.Model(&models.Todo{}).Where("list_id = ?", id).Update("list_id", 0)
	return database.DB.Delete(&models.List{}, id).Error
}
