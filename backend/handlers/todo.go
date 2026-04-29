package handlers

import (
	"net/http"
	"strconv"
	"todo-list/backend/models"
	"todo-list/backend/services"

	"github.com/gin-gonic/gin"
)

func parseListID(c *gin.Context) uint {
	idStr := c.Query("list_id")
	if idStr == "" {
		return 0
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0
	}
	return uint(id)
}

func ListTodos(c *gin.Context) {
	listID := parseListID(c)
	status := c.Query("status")
	priority := c.Query("priority")
	tag := c.Query("tag")
	search := c.Query("search")

	todos, err := services.GetTodos(listID, status, priority, tag, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func GetTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	todo, err := services.GetTodo(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateTodo(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := services.UpdateTodo(uint(id), &input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func ToggleTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	todo, err := services.ToggleTodo(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := services.DeleteTodo(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

type reorderRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

func ReorderTodos(c *gin.Context) {
	var req reorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ReorderTodos(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "reordered"})
}

func ExportTodos(c *gin.Context) {
	format := c.DefaultQuery("format", "json")
	listID := parseListID(c)

	todos, err := services.GetTodos(listID, "", "", "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if format == "csv" {
		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=todos.csv")
		c.String(http.StatusOK, "ID,标题,备注,优先级,标签,已完成,截止日期,创建日期\n")
		for _, t := range todos {
			dueDate := ""
			if t.DueDate.Valid {
				dueDate = t.DueDate.Time.Format("2006-01-02")
			}
			c.String(http.StatusOK, "%d,\"%s\",\"%s\",%s,%s,%v,%s,%s\n",
				t.ID, t.Title, t.Description, t.Priority, t.Tags, t.Completed,
				dueDate, t.CreatedAt.Format("2006-01-02"))
		}
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=todos.json")
	c.JSON(http.StatusOK, todos)
}

func Stats(c *gin.Context) {
	listID := parseListID(c)
	stats := services.GetStats(listID)
	c.JSON(http.StatusOK, stats)
}
