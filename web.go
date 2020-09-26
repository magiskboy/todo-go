package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
)

// App HTTP application
var App *gin.Engine
var identifyKey = "email"

// CheckHealth api check healthy
func CheckHealth(ctx *gin.Context) {
	ctx.Writer.WriteHeader(http.StatusNoContent)
}

// CreateUserRequest json data for create a new user
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateUser user login handler
func CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    400,
		})
		return
	}
	newUser, err := CreateNewUser(req.Email, req.Name, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    400,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":   newUser.ID,
		"code": 200,
	})
}

// GetTasks get all tasks of a user
func GetTasks(ctx *gin.Context) {
	email, _ := ctx.Get(identifyKey)
	user, _ := GetUserByEmail(email.(string))
	ctx.JSON(http.StatusOK, gin.H{
		"code":  200,
		"tasks": user.Tasks,
	})
}

// CreateTaskRequest use to create task
type CreateTaskRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateTask create a new task
func CreateTask(ctx *gin.Context) {
	var req CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	email, _ := ctx.Get(identifyKey)
	user, _ := GetUserByEmail(email.(string))
	task, err := user.AddTask(req.Name, req.Description, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"task_id": task.ID,
	})
}

// RemoveTask delete a task of the user
func RemoveTask(ctx *gin.Context) {
	taskID, err := strconv.Atoi(ctx.Param("task_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Task ID isn't valid",
			"code":    400,
		})
		return
	}
	TaskID := uint(taskID)
	email, _ := ctx.Get(identifyKey)
	user, _ := GetUserByEmail(email.(string))
	_, err = user.RemoveTask(TaskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"code":    500,
		})
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
}

// UpdateTaskRequest json
type UpdateTaskRequest struct {
	CreateTaskRequest
	Done bool `json:"done"`
}

// UpdateTask update a task of the user
func UpdateTask(ctx *gin.Context) {
	var req UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	TaskID := ctx.Param("task_id")
	task, err := GetTaskByID(TaskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	UserID, _ := ctx.Get("ID")
	if task.UserRefer != UserID.(uint) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    401,
			"message": "Task owned by other",
		})
		return
	}
	task.Name = req.Name
	task.Description = req.Description
	task.Done = req.Done
	db.Save(&task)
	ctx.Writer.WriteHeader(http.StatusOK)
}

// InitHTTP initialize function
func InitHTTP() {
	authMiddleware, err := CreateAuthMiddleware(os.Getenv("SECRET_KEY"), identifyKey)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(-1)
	}

	App = gin.Default()
	App.GET("/health", CheckHealth)
	public := App.Group("/public")
	{
		public.POST("/login", authMiddleware.LoginHandler)
		public.POST("/users", CreateUser)
	}
	private := App.Group("/private")
	private.Use(authMiddleware.MiddlewareFunc())
	{
		private.GET("/tasks", GetTasks)
		private.POST("/tasks", CreateTask)
		private.PUT("/tasks/:task_id", UpdateTask)
		private.DELETE("/tasks/:task_id", RemoveTask)
	}
}
