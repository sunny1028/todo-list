package services

import (
	"errors"
	"todo-list/backend/database"
	"todo-list/backend/models"
	"todo-list/backend/repository"
)

func CreateShare(userID uint, listID uint, permission string) (*models.ListShare, error) {
	list, err := repository.FindListByID(userID, listID)
	if err != nil {
		return nil, errors.New("list not found")
	}
	if list.UserID != userID {
		return nil, errors.New("only the list owner can share")
	}
	if permission != "view" && permission != "edit" {
		permission = "view"
	}
	return repository.CreateShare(listID, permission)
}

func GetShare(userID uint, listID uint) (*models.ListShare, []models.ListShareMember, error) {
	list, err := repository.FindListByID(userID, listID)
	if err != nil {
		return nil, nil, errors.New("list not found")
	}
	if list.UserID != userID {
		return nil, nil, errors.New("only the list owner can view share info")
	}
	s, err := repository.FindShareByListID(listID)
	if err != nil {
		return nil, nil, errors.New("share not found")
	}
	members, _ := repository.FindShareMembers(s.ID)
	return s, members, nil
}

func DeleteShare(userID uint, listID uint) error {
	list, err := repository.FindListByID(userID, listID)
	if err != nil {
		return errors.New("list not found")
	}
	if list.UserID != userID {
		return errors.New("only the list owner can revoke share")
	}
	return repository.DeleteShare(listID)
}

func JoinShare(userID uint, code string) (*models.List, error) {
	s, err := repository.FindShareByCode(code)
	if err != nil {
		return nil, errors.New("invalid share code")
	}
	var list models.List
	if database.DB.First(&list, s.ListID).Error != nil {
		return nil, errors.New("list not found")
	}
	if list.UserID == userID {
		return nil, errors.New("cannot join your own list")
	}
	if err := repository.AddShareMember(s.ID, userID); err != nil {
		return nil, err
	}
	return &list, nil
}
