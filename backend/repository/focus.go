package repository

import (
	"time"
	"todo-list/backend/database"
	"todo-list/backend/models"
)

func CreateFocusSession(s *models.FocusSession) error {
	return database.DB.Create(s).Error
}

func UpdateFocusSession(s *models.FocusSession) error {
	return database.DB.Save(s).Error
}

func FindFocusSessions(userID uint, limit int) ([]models.FocusSession, error) {
	var sessions []models.FocusSession
	err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	for i := range sessions {
		if sessions[i].TodoID != nil {
			var todo models.Todo
			if database.DB.Select("title").First(&todo, *sessions[i].TodoID).Error == nil {
				sessions[i].TodoTitle = todo.Title
			}
		}
	}
	return sessions, nil
}

type FocusStats struct {
	TodayMinutes  int   `json:"today_minutes"`
	TotalMinutes  int   `json:"total_minutes"`
	TotalSessions int64 `json:"total_sessions"`
	StreakDays    int   `json:"streak_days"`
}

func GetFocusStats(userID uint) FocusStats {
	var s FocusStats
	today := time.Now().Format("2006-01-02")

	database.DB.Model(&models.FocusSession{}).
		Where("user_id = ? AND completed = ? AND date(started_at) = ?", userID, true, today).
		Select("COALESCE(SUM(duration_min), 0)").Scan(&s.TodayMinutes)

	database.DB.Model(&models.FocusSession{}).
		Where("user_id = ? AND completed = ?", userID, true).
		Select("COALESCE(SUM(duration_min), 0)").Scan(&s.TotalMinutes)

	database.DB.Model(&models.FocusSession{}).
		Where("user_id = ? AND completed = ?", userID, true).Count(&s.TotalSessions)

	d := time.Now()
	for {
		dateStr := d.Format("2006-01-02")
		var count int64
		database.DB.Model(&models.FocusSession{}).
			Where("user_id = ? AND completed = ? AND date(started_at) = ?", userID, true, dateStr).
			Count(&count)
		if count == 0 {
			break
		}
		s.StreakDays++
		d = d.AddDate(0, 0, -1)
	}

	return s
}
