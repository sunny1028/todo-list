package services

import (
	"errors"
	"todo-list/backend/models"
	"todo-list/backend/repository"
)

func GetTodos(userID uint, listID uint, status, priority, tag, search string) ([]models.Todo, error) {
	return repository.FindAll(userID, listID, status, priority, tag, search)
}

func GetStats(userID uint, listID uint) map[string]interface{} {
	total, active, completed := repository.CountAll(userID, listID)
	byPriority := repository.CountByPriority(userID, listID)
	byTag := repository.CountByTag(userID, listID)

	return map[string]interface{}{
		"total":       total,
		"active":      active,
		"completed":   completed,
		"by_priority": byPriority,
		"by_tag":      byTag,
	}
}

func GetTodo(userID uint, id uint) (*models.Todo, error) {
	todo, err := repository.FindByID(userID, id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

func CreateTodo(todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("title is required")
	}
	if todo.Priority == "" {
		todo.Priority = "medium"
	}
	maxOrder, _ := repository.GetMaxSortOrder(todo.UserID)
	todo.SortOrder = maxOrder + 1
	return repository.Create(todo)
}

func UpdateTodo(userID uint, id uint, input *models.Todo) (*models.Todo, error) {
	todo, err := repository.FindByID(userID, id)
	if err != nil {
		return nil, errors.New("todo not found")
	}

	todo.Title = input.Title
	todo.Description = input.Description
	todo.Priority = input.Priority
	todo.Tags = input.Tags
	todo.DueDate = input.DueDate

	if err := repository.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func ToggleTodo(userID uint, id uint) (*models.Todo, error) {
	todo, err := repository.FindByID(userID, id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	todo.Completed = !todo.Completed
	if err := repository.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func ReorderTodos(userID uint, ids []uint) error {
	for i, id := range ids {
		if err := repository.UpdateOrder(userID, id, i); err != nil {
			return err
		}
	}
	return nil
}

func ArchiveTodo(userID uint, id uint) (*models.Todo, error) {
	todo, err := repository.FindByID(userID, id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	todo.Archived = true
	if err := repository.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func UnarchiveTodo(userID uint, id uint) (*models.Todo, error) {
	todo, err := repository.FindByID(userID, id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	todo.Archived = false
	if err := repository.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func DeleteTodo(userID uint, id uint) error {
	return repository.Delete(userID, id)
}
