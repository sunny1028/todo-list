package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"todo-list/backend/database"
	"todo-list/backend/models"

	"github.com/gin-gonic/gin"
)

func UploadAttachment(c *gin.Context) {
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	uploadDir := "uploads"
	os.MkdirAll(uploadDir, 0755)

	ext := filepath.Ext(file.Filename)
	savedName := fmt.Sprintf("%d_%d%s", todoID, time.Now().UnixNano(), ext)
	savedPath := filepath.Join(uploadDir, savedName)

	if err := c.SaveUploadedFile(file, savedPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	att := models.Attachment{
		TodoID:   uint(todoID),
		Filename: file.Filename,
		Filepath: savedPath,
		MimeType: file.Header.Get("Content-Type"),
		Size:     file.Size,
	}
	database.DB.Create(&att)

	c.JSON(http.StatusCreated, att)
}

func ListAttachments(c *gin.Context) {
	todoID := c.Param("id")
	var atts []models.Attachment
	database.DB.Where("todo_id = ?", todoID).Find(&atts)
	c.JSON(http.StatusOK, atts)
}

func ServeAttachment(c *gin.Context) {
	id := c.Param("id")
	var att models.Attachment
	if err := database.DB.First(&att, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.File(att.Filepath)
}

func DeleteAttachment(c *gin.Context) {
	id := c.Param("id")
	var att models.Attachment
	if err := database.DB.First(&att, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	os.Remove(att.Filepath)
	database.DB.Delete(&att)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
