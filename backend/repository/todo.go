package repository

import (
	"strings"
	"time"
	"todo-list/backend/database"
	"todo-list/backend/models"

	"gorm.io/gorm"
)

func FindAllByList(listID uint, status, priority, tag, search string) ([]models.Todo, error) {
	var todos []models.Todo
	q := database.DB.Model(&models.Todo{}).Where("list_id = ?", listID)

	if status == "archived" {
		q = q.Where("archived = ?", true)
	} else {
		q = q.Where("archived = ?", false)
		if status == "completed" {
			q = q.Where("completed = ?", true)
		} else if status == "active" {
			q = q.Where("completed = ?", false)
		}
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
	if err != nil {
		return todos, err
	}

	if len(todos) > 0 {
		ids := make([]uint, len(todos))
		for i, t := range todos {
			ids[i] = t.ID
		}
		type countRow struct {
			TodoID    uint
			Total     int
			Completed int
		}
		var rows []countRow
		database.DB.Model(&models.Subtask{}).
			Select("todo_id, count(*) as total, sum(case when completed then 1 else 0 end) as completed").
			Where("todo_id IN ?", ids).
			Group("todo_id").
			Scan(&rows)
		for i := range todos {
			for _, r := range rows {
				if r.TodoID == todos[i].ID {
					todos[i].SubtaskCount = r.Total
					todos[i].SubtaskCompleted = r.Completed
					break
				}
			}
		}
	}

	return todos, nil
}

func FindAll(userID uint, listID uint, status, priority, tag, search string) ([]models.Todo, error) {
	var todos []models.Todo
	q := database.DB.Model(&models.Todo{}).Where("user_id = ?", userID)

	if status == "archived" {
		q = q.Where("archived = ?", true)
	} else {
		q = q.Where("archived = ?", false)
		if status == "completed" {
			q = q.Where("completed = ?", true)
		} else if status == "active" {
			q = q.Where("completed = ?", false)
		}
	}

	if listID > 0 {
		q = q.Where("list_id = ?", listID)
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
	if err != nil {
		return todos, err
	}

	if len(todos) > 0 {
		ids := make([]uint, len(todos))
		for i, t := range todos {
			ids[i] = t.ID
		}
		type countRow struct {
			TodoID    uint
			Total     int
			Completed int
		}
		var rows []countRow
		database.DB.Model(&models.Subtask{}).
			Select("todo_id, count(*) as total, sum(case when completed then 1 else 0 end) as completed").
			Where("todo_id IN ?", ids).
			Group("todo_id").
			Scan(&rows)
		for i := range todos {
			for _, r := range rows {
				if r.TodoID == todos[i].ID {
					todos[i].SubtaskCount = r.Total
					todos[i].SubtaskCompleted = r.Completed
					break
				}
			}
		}
	}

	return todos, nil
}

func FindByID(userID uint, id uint) (*models.Todo, error) {
	var todo models.Todo
	err := database.DB.Where("user_id = ?", userID).First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	// Populate subtask counts
	var total, completed int64
	database.DB.Model(&models.Subtask{}).Where("todo_id = ?", id).Count(&total)
	database.DB.Model(&models.Subtask{}).Where("todo_id = ? AND completed = ?", id, true).Count(&completed)
	todo.SubtaskCount = int(total)
	todo.SubtaskCompleted = int(completed)
	return &todo, nil
}

func Create(todo *models.Todo) error {
	return database.DB.Create(todo).Error
}

func Update(todo *models.Todo) error {
	return database.DB.Save(todo).Error
}

func UpdateOrder(userID uint, id uint, sortOrder int) error {
	return database.DB.Model(&models.Todo{}).Where("user_id = ? AND id = ?", userID, id).Update("sort_order", sortOrder).Error
}

func Delete(userID uint, id uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&models.Todo{}, id).Error
}

func GetMaxSortOrder(userID uint) (int, error) {
	var maxOrder int
	err := database.DB.Model(&models.Todo{}).Where("user_id = ?", userID).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder).Error
	return maxOrder, err
}

func CountAll(userID uint, listID uint) (total, active, completed int64) {
	q := database.DB.Model(&models.Todo{}).Where("user_id = ?", userID)
	if listID > 0 {
		q = q.Where("list_id = ?", listID)
	}
	q.Count(&total)
	q.Where("completed = ?", false).Count(&active)
	database.DB.Model(&models.Todo{}).Where("user_id = ? AND completed = ?", userID, true).Count(&completed)
	return
}

func CountByPriority(userID uint, listID uint) map[string]int64 {
	var results []struct {
		Priority string
		Count    int64
	}
	q := database.DB.Model(&models.Todo{}).Where("user_id = ?", userID).Select("priority, count(*) as count").Group("priority")
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

type DailyTrend struct {
	Date      string `json:"date"`
	Created   int64  `json:"created"`
	Completed int64  `json:"completed"`
}

func GetDailyTrends(userID uint, days int) []DailyTrend {
	var results []DailyTrend
	for i := 0; i < days; i++ {
		d := time.Now().AddDate(0, 0, -i)
		dateStr := d.Format("2006-01-02")
		var created, completed int64
		database.DB.Model(&models.Todo{}).Where("user_id = ? AND date(created_at) = ?", userID, dateStr).Count(&created)
		database.DB.Model(&models.Todo{}).Where("user_id = ? AND completed = ? AND date(updated_at) = ?", userID, true, dateStr).Count(&completed)
		results = append(results, DailyTrend{Date: dateStr, Created: created, Completed: completed})
	}
	// Reverse so oldest first
	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
	}
	return results
}

func buildQuery(userID uint, listID uint) *gorm.DB {
	q := database.DB.Model(&models.Todo{}).Where("user_id = ?", userID)
	if listID > 0 {
		q = q.Where("list_id = ?", listID)
	}
	return q
}

func CountCreatedOnDate(userID uint, listID uint, date string) int64 {
	var c int64
	buildQuery(userID, listID).Where("date(created_at) = ?", date).Count(&c)
	return c
}

func CountCompletedOnDate(userID uint, listID uint, date string) int64 {
	var c int64
	buildQuery(userID, listID).Where("completed = ? AND date(updated_at) = ?", true, date).Count(&c)
	return c
}

func CountCreatedCompletedOnDate(userID uint, listID uint, date string) (created, completed int64) {
	buildQuery(userID, listID).Where("date(created_at) = ?", date).Count(&created)
	buildQuery(userID, listID).Where("completed = ? AND date(updated_at) = ?", true, date).Count(&completed)
	return
}

func CountByTag(userID uint, listID uint) map[string]int64 {
	var todos []models.Todo
	q := database.DB.Model(&models.Todo{}).Where("user_id = ? AND tags != ''", userID)
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
