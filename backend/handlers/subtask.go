package handlers

import (
	"net/http"
	"strconv"

	"todo-list/backend/database"
	"todo-list/backend/models"

	"github.com/gin-gonic/gin"
)

func ListSubtasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	todoID := c.Param("id")
	var subtasks []models.Subtask
	database.DB.Where("user_id = ? AND todo_id = ?", userID, todoID).Find(&subtasks)
	c.JSON(http.StatusOK, subtasks)
}

func CreateSubtask(c *gin.Context) {
	userID := c.GetUint("user_id")
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
	st.UserID = userID
	st.TodoID = uint(todoID)
	database.DB.Create(&st)
	c.JSON(http.StatusCreated, st)
}

func ToggleSubtask(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var st models.Subtask
	if err := database.DB.Where("user_id = ?", userID).First(&st, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	st.Completed = !st.Completed
	database.DB.Save(&st)

	// Auto-complete parent todo when all subtasks are done
	if st.Completed {
		var count int64
		database.DB.Model(&models.Subtask{}).Where("todo_id = ? AND completed = ?", st.TodoID, false).Count(&count)
		if count == 0 {
			database.DB.Model(&models.Todo{}).Where("id = ?", st.TodoID).Update("completed", true)
		}
	} else {
		// Uncheck parent if subtask is unchecked
		database.DB.Model(&models.Todo{}).Where("id = ?", st.TodoID).Update("completed", false)
	}

	c.JSON(http.StatusOK, st)
}

func UpdateSubtask(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var st models.Subtask
	if err := database.DB.Where("user_id = ?", userID).First(&st, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var input models.Subtask
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Title != "" {
		st.Title = input.Title
	}
	database.DB.Save(&st)
	c.JSON(http.StatusOK, st)
}

func DeleteSubtask(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	database.DB.Where("user_id = ?", userID).Delete(&models.Subtask{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
