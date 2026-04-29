package services

import (
	"errors"
	"todo-list/backend/models"
	"todo-list/backend/repository"
)

func GetTodos(listID uint, status, priority, tag, search string) ([]models.Todo, error) {
	return repository.FindAll(listID, status, priority, tag, search)
}

func GetStats(listID uint) map[string]interface{} {
	total, active, completed := repository.CountAll(listID)
	byPriority := repository.CountByPriority(listID)
	byTag := repository.CountByTag(listID)

	return map[string]interface{}{
		"total":       total,
		"active":      active,
		"completed":   completed,
		"by_priority": byPriority,
		"by_tag":      byTag,
	}
}

func GetTodo(id uint) (*models.Todo, error) {
	todo, err := repository.FindByID(id)
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
	maxOrder, _ := repository.GetMaxSortOrder()
	todo.SortOrder = maxOrder + 1
	return repository.Create(todo)
}

func UpdateTodo(id uint, input *models.Todo) (*models.Todo, error) {
	todo, err := repository.FindByID(id)
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

func ToggleTodo(id uint) (*models.Todo, error) {
	todo, err := repository.FindByID(id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	todo.Completed = !todo.Completed
	if err := repository.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func ReorderTodos(ids []uint) error {
	for i, id := range ids {
		if err := repository.UpdateOrder(id, i); err != nil {
			return err
		}
	}
	return nil
}

func ArchiveTodo(id uint) (*models.Todo, error) {
	todo, err := repository.FindByID(id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	todo.Archived = true
	if err := repository.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func UnarchiveTodo(id uint) (*models.Todo, error) {
	todo, err := repository.FindByID(id)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	todo.Archived = false
	if err := repository.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func DeleteTodo(id uint) error {
	return repository.Delete(id)
}
