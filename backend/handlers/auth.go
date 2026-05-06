package handlers

import (
	"net/http"
	"time"

	"todo-list/backend/config"
	"todo-list/backend/database"
	"todo-list/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func generateJWT(userID uint, uuid string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"uuid":    uuid,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Load().JWTSecret))
}

type uuidRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

func AuthUUID(c *gin.Context) {
	var req uuidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid required"})
		return
	}

	var user models.User
	result := database.DB.Where("uuid = ?", req.UUID).First(&user)

	if result.Error != nil {
		// New user — check if this is the first user to migrate legacy data
		var count int64
		database.DB.Model(&models.User{}).Count(&count)
		isFirst := count == 0

		user = models.User{UUID: req.UUID}
		database.DB.Create(&user)

		// First user claims all user_id=0 orphaned data
		if isFirst {
			uid := user.ID
			database.DB.Model(&models.Todo{}).Where("user_id = 0").Update("user_id", uid)
			database.DB.Model(&models.List{}).Where("user_id = 0").Update("user_id", uid)
			database.DB.Model(&models.Subtask{}).Where("user_id = 0").Update("user_id", uid)
			database.DB.Model(&models.Attachment{}).Where("user_id = 0").Update("user_id", uid)
		}
	}

	token, err := generateJWT(user.ID, user.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":        token,
		"has_password": user.HasPassword,
	})
}

type bindRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AuthBind(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req bindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	if len(req.Password) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 4 characters"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// Check username not taken by another user
	var existing models.User
	if database.DB.Where("username = ? AND id != ?", req.Username, userID).First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	}

	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"username":      req.Username,
		"password_hash": string(hash),
		"has_password":  true,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bind failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AuthLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	var user models.User
	if database.DB.Where("username = ?", req.Username).First(&user).Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := generateJWT(user.ID, user.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":        token,
		"has_password": user.HasPassword,
	})
}
