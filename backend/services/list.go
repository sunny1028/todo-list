package services

import (
	"errors"
	"todo-list/backend/models"
	"todo-list/backend/repository"
)

func GetLists(userID uint) ([]models.List, error) {
	own, err := repository.FindAllLists(userID)
	if err != nil {
		return nil, err
	}
	sharedIDs := repository.FindSharedListIDs(userID)
	if len(sharedIDs) > 0 {
		shared, err := repository.FindListsByIDs(sharedIDs)
		if err == nil {
			for i := range shared {
				shared[i].Permission = repository.GetUserPermission(shared[i].ID, userID)
			}
			own = append(own, shared...)
		}
	}
	return own, nil
}

func GetList(userID uint, id uint) (*models.List, error) {
	list, err := repository.FindListByID(userID, id)
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

func UpdateList(userID uint, id uint, input *models.List) (*models.List, error) {
	list, err := repository.FindListByID(userID, id)
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

func DeleteList(userID uint, id uint) error {
	return repository.DeleteList(userID, id)
}
