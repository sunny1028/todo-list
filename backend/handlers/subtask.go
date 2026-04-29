package handlers

import (
	"net/http"
	"strconv"

	"todo-list/backend/database"
	"todo-list/backend/models"

	"github.com/gin-gonic/gin"
)

func ListSubtasks(c *gin.Context) {
	todoID := c.Param("id")
	var subtasks []models.Subtask
	database.DB.Where("todo_id = ?", todoID).Find(&subtasks)
	c.JSON(http.StatusOK, subtasks)
}

func CreateSubtask(c *gin.Context) {
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}
	var st models.Subtask
	if err := c.ShouldBindJSON(&st); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	st.TodoID = uint(todoID)
	database.DB.Create(&st)
	c.JSON(http.StatusCreated, st)
}

func ToggleSubtask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var st models.Subtask
	if err := database.DB.First(&st, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	st.Completed = !st.Completed
	database.DB.Save(&st)
	c.JSON(http.StatusOK, st)
}

func DeleteSubtask(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Subtask{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
