package handlers

import (
	"net/http"
	"strconv"
	"todo-list/backend/services"

	"github.com/gin-gonic/gin"
)

type startFocusRequest struct {
	TodoID      *uint `json:"todo_id"`
	DurationMin int   `json:"duration_min" binding:"required"`
}

func StartFocus(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req startFocusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "duration_min is required"})
		return
	}
	s, err := services.StartFocus(userID, req.TodoID, req.DurationMin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func CompleteFocus(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	s, err := services.CompleteFocus(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func ListFocusSessions(c *gin.Context) {
	userID := c.GetUint("user_id")
	sessions, err := services.GetFocusSessions(userID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func GetFocusStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	stats := services.GetFocusStats(userID)
	c.JSON(http.StatusOK, stats)
}
