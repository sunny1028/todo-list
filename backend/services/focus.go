package services

import (
	"errors"
	"time"
	"todo-list/backend/database"
	"todo-list/backend/models"
	"todo-list/backend/repository"

	"gorm.io/gorm"
)

func StartFocus(userID uint, todoID *uint, durationMin int) (*models.FocusSession, error) {
	if durationMin <= 0 || durationMin > 180 {
		return nil, errors.New("duration must be 1-180 minutes")
	}
	s := &models.FocusSession{
		UserID:      userID,
		TodoID:      todoID,
		DurationMin: durationMin,
		StartedAt:   time.Now(),
	}
	if err := repository.CreateFocusSession(s); err != nil {
		return nil, err
	}
	return s, nil
}

func CompleteFocus(userID uint, id uint) (*models.FocusSession, error) {
	var s models.FocusSession
	if err := database.DB.Where("user_id = ?", userID).First(&s, id).Error; err != nil {
		return nil, errors.New("session not found")
	}
	if s.Completed {
		return &s, nil
	}
	now := time.Now()
	s.Completed = true
	s.EndedAt = &now
	if err := repository.UpdateFocusSession(&s); err != nil {
		return nil, err
	}
	if s.TodoID != nil {
		database.DB.Model(&models.Todo{}).Where("id = ?", *s.TodoID).
			Update("focus_minutes", gorm.Expr("focus_minutes + ?", s.DurationMin))
	}
	return &s, nil
}

func GetFocusSessions(userID uint, limit int) ([]models.FocusSession, error) {
	return repository.FindFocusSessions(userID, limit)
}

func GetFocusStats(userID uint) repository.FocusStats {
	return repository.GetFocusStats(userID)
}
