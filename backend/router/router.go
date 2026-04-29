package router

import (
	"todo-list/backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(corsOrigin string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/todos", handlers.ListTodos)
		api.POST("/todos", handlers.CreateTodo)
		api.GET("/todos/stats", handlers.Stats)
		api.GET("/todos/export", handlers.ExportTodos)
		api.PUT("/todos/reorder", handlers.ReorderTodos)
		api.GET("/todos/:id", handlers.GetTodo)
		api.PUT("/todos/:id", handlers.UpdateTodo)
		api.PATCH("/todos/:id/toggle", handlers.ToggleTodo)
		api.DELETE("/todos/:id", handlers.DeleteTodo)

		api.GET("/lists", handlers.ListLists)
		api.POST("/lists", handlers.CreateList)
		api.PUT("/lists/:id", handlers.UpdateList)
		api.DELETE("/lists/:id", handlers.DeleteList)

		api.GET("/todos/:id/attachments", handlers.ListAttachments)
		api.POST("/todos/:id/attachments", handlers.UploadAttachment)
		api.GET("/attachments/:id", handlers.ServeAttachment)
		api.DELETE("/attachments/:id", handlers.DeleteAttachment)
	}

	return r
}
