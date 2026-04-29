package services

import (
	"errors"
	"todo-list/backend/models"
	"todo-list/backend/repository"
)

func GetLists() ([]models.List, error) {
	return repository.FindAllLists()
}

func GetList(id uint) (*models.List, error) {
	list, err := repository.FindListByID(id)
	if err != nil {
		return nil, errors.New("list not found")
	}
	return list, nil
}

func CreateList(list *models.List) error {
	if list.Name == "" {
		return errors.New("name is required")
	}
	return repository.CreateList(list)
}

func UpdateList(id uint, input *models.List) (*models.List, error) {
	list, err := repository.FindListByID(id)
	if err != nil {
		return nil, errors.New("list not found")
	}
	list.Name = input.Name
	list.Color = input.Color
	if err := repository.UpdateList(list); err != nil {
		return nil, err
	}
	return list, nil
}

func DeleteList(id uint) error {
	return repository.DeleteList(id)
}
