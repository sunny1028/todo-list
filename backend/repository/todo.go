package repository

import (
	"strings"
	"todo-list/backend/database"
	"todo-list/backend/models"
)

func FindAll(listID uint, status, priority, tag, search string) ([]models.Todo, error) {
	var todos []models.Todo
	q := database.DB.Model(&models.Todo{})

	if listID > 0 {
		q = q.Where("list_id = ?", listID)
	}

	if status == "completed" {
		q = q.Where("completed = ?", true)
	} else if status == "active" {
		q = q.Where("completed = ?", false)
	}

	if priority != "" {
		q = q.Where("priority = ?", priority)
	}

	if tag != "" {
		q = q.Where("tags LIKE ?", "%"+tag+"%")
	}

	if search != "" {
		q = q.Where("title LIKE ?", "%"+search+"%")
	}

	err := q.Order("sort_order ASC, created_at DESC").Find(&todos).Error
	return todos, err
}

func FindByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := database.DB.First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func Create(todo *models.Todo) error {
	return database.DB.Create(todo).Error
}

func Update(todo *models.Todo) error {
	return database.DB.Save(todo).Error
}

func UpdateOrder(id uint, sortOrder int) error {
	return database.DB.Model(&models.Todo{}).Where("id = ?", id).Update("sort_order", sortOrder).Error
}

func Delete(id uint) error {
	return database.DB.Delete(&models.Todo{}, id).Error
}

func GetMaxSortOrder() (int, error) {
	var maxOrder int
	err := database.DB.Model(&models.Todo{}).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder).Error
	return maxOrder, err
}

// Stats helpers
func CountAll(listID uint) (total, active, completed int64) {
	q := database.DB.Model(&models.Todo{})
	if listID > 0 {
		q = q.Where("list_id = ?", listID)
	}
	q.Count(&total)
	q.Where("completed = ?", false).Count(&active)
	q.Where("completed = ?", true).Count(&completed)
	return
}

func CountByPriority(listID uint) map[string]int64 {
	var results []struct {
		Priority string
		Count    int64
	}
	q := database.DB.Model(&models.Todo{}).Select("priority, count(*) as count").Group("priority")
	if listID > 0 {
		q = q.Where("list_id = ?", listID)
	}
	q.Scan(&results)
	out := make(map[string]int64)
	for _, r := range results {
		out[r.Priority] = r.Count
	}
	return out
}

func CountByTag(listID uint) map[string]int64 {
	var todos []models.Todo
	q := database.DB.Model(&models.Todo{}).Where("tags != ''")
	if listID > 0 {
		q = q.Where("list_id = ?", listID)
	}
	q.Find(&todos)
	counts := make(map[string]int64)
	for _, t := range todos {
		for _, tag := range strings.Split(t.Tags, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				counts[tag]++
			}
		}
	}
	return counts
}
