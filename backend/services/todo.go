package services

import (
	"errors"
	"time"
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
	trends := repository.GetDailyTrends(userID, 7)

	return map[string]interface{}{
		"total":       total,
		"active":      active,
		"completed":   completed,
		"by_priority": byPriority,
		"by_tag":      byTag,
		"daily_trends": trends,
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
	todo.Effort = input.Effort
	todo.Tags = input.Tags
	todo.DueDate = input.DueDate
	todo.Recurrence = input.Recurrence

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

	// If completing a recurring task, spawn next instance
	if todo.Completed && todo.Recurrence != "" {
		next := models.Todo{
			UserID:      todo.UserID,
			ListID:      todo.ListID,
			Title:       todo.Title,
			Description: todo.Description,
			Priority:    todo.Priority,
			Effort:      todo.Effort,
			Tags:        todo.Tags,
			Recurrence:  todo.Recurrence,
		}
		d := todo.DueDate.Time
		if !todo.DueDate.Valid {
			d = time.Now()
		}
		switch todo.Recurrence {
		case "daily":
			next.DueDate = models.DateOnly{Time: d.AddDate(0, 0, 1), Valid: true}
		case "weekly":
			next.DueDate = models.DateOnly{Time: d.AddDate(0, 0, 7), Valid: true}
		case "monthly":
			next.DueDate = models.DateOnly{Time: d.AddDate(0, 1, 0), Valid: true}
		}
		CreateTodo(&next)
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
