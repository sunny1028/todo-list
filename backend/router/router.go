package router

import (
	"todo-list/backend/handlers"
	"todo-list/backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(corsOrigin string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Public auth routes
	api := r.Group("/api")
	{
		api.POST("/auth/uuid", handlers.AuthUUID)
		api.POST("/auth/login", handlers.AuthLogin)
	}

	// Protected routes (JWT required)
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		protected.POST("/auth/bind", handlers.AuthBind)
		protected.POST("/auth/merge", handlers.AuthMerge)

		// Todos
		protected.GET("/todos", handlers.ListTodos)
		protected.POST("/todos", handlers.CreateTodo)
		protected.GET("/todos/stats", handlers.Stats)
		protected.GET("/todos/export", handlers.ExportTodos)
		protected.POST("/todos/import", handlers.ImportTodos)
		protected.PUT("/todos/reorder", handlers.ReorderTodos)
		protected.GET("/todos/:id", handlers.GetTodo)
		protected.PUT("/todos/:id", handlers.UpdateTodo)
		protected.PATCH("/todos/:id/toggle", handlers.ToggleTodo)
		protected.PATCH("/todos/:id/archive", handlers.ArchiveTodo)
		protected.PATCH("/todos/:id/unarchive", handlers.UnarchiveTodo)
		protected.DELETE("/todos/:id", handlers.DeleteTodo)

		// Subtasks
		protected.GET("/todos/:id/subtasks", handlers.ListSubtasks)
		protected.POST("/todos/:id/subtasks", handlers.CreateSubtask)
		protected.PATCH("/subtasks/:id", handlers.UpdateSubtask)
		protected.PATCH("/subtasks/:id/toggle", handlers.ToggleSubtask)
		protected.PUT("/subtasks/reorder", handlers.ReorderSubtasks)
		protected.DELETE("/subtasks/:id", handlers.DeleteSubtask)

		// Attachments
		protected.GET("/todos/:id/attachments", handlers.ListAttachments)
		protected.POST("/todos/:id/attachments", handlers.UploadAttachment)
		protected.GET("/attachments/:id", handlers.ServeAttachment)
		protected.DELETE("/attachments/:id", handlers.DeleteAttachment)

		// Lists
		protected.GET("/lists", handlers.ListLists)
		protected.POST("/lists", handlers.CreateList)
		protected.PUT("/lists/:id", handlers.UpdateList)
		protected.DELETE("/lists/:id", handlers.DeleteList)
	}

	return r
}
