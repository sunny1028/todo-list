package handlers

import (
	"net/http"
	"strconv"
	"todo-list/backend/services"

	"github.com/gin-gonic/gin"
)

func CreateShare(c *gin.Context) {
	userID := c.GetUint("user_id")
	listID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Permission string `json:"permission"`
	}
	c.ShouldBindJSON(&body)
	share, err := services.CreateShare(userID, uint(listID), body.Permission)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, share)
}

func GetShare(c *gin.Context) {
	userID := c.GetUint("user_id")
	listID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	share, members, err := services.GetShare(userID, uint(listID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"share":   share,
		"members": members,
	})
}

func DeleteShare(c *gin.Context) {
	userID := c.GetUint("user_id")
	listID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := services.DeleteShare(userID, uint(listID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "share revoked"})
}

func JoinList(c *gin.Context) {
	userID := c.GetUint("user_id")
	var body struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}
	list, err := services.JoinShare(userID, body.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
