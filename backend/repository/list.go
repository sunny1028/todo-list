package repository

import (
	"todo-list/backend/database"
	"todo-list/backend/models"
)

func FindAllLists(userID uint) ([]models.List, error) {
	var lists []models.List
	err := database.DB.Where("user_id = ?", userID).Order("created_at ASC").Find(&lists).Error
	return lists, err
}

func FindListByID(userID uint, id uint) (*models.List, error) {
	var list models.List
	err := database.DB.Where("user_id = ?", userID).First(&list, id).Error
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

func FindListsByIDs(ids []uint) ([]models.List, error) {
	var lists []models.List
	err := database.DB.Where("id IN ?", ids).Find(&lists).Error
	return lists, err
}

func DeleteList(userID uint, id uint) error {
	database.DB.Model(&models.Todo{}).Where("user_id = ? AND list_id = ?", userID, id).Update("list_id", 0)
	return database.DB.Where("user_id = ?", userID).Delete(&models.List{}, id).Error
}
