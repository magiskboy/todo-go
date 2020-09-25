package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

// App is a HTTP application
var App *gin.Engine

var identifyKey = "email"

// LoginInfo json login payload
type LoginInfo struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserInfo json data for create a new user
type UserInfo struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TaskInfo use to create task
type TaskInfo struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CheckHealth api check healthy
func CheckHealth(ctx *gin.Context) {
	ctx.Writer.WriteHeader(http.StatusNoContent)
}

// CreateUser user login handler
func CreateUser(ctx *gin.Context) {
	var data UserInfo
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"code":    400,
		})
		return
	}
	newUser, err := CreateNewUser(data.Email, data.Name, data.Password)
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
	user, _ := ctx.Get(identifyKey)
	tasks := user.(*User).Tasks
	if tasks == nil {
		tasks = make([]Task, 0)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":  200,
		"tasks": tasks,
	})
}

// CreateTask create a new task
func CreateTask(ctx *gin.Context) {
	var data TaskInfo
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	user, _ := ctx.Get(identifyKey)
	task, err := user.(*User).AddTask(data.Name, data.Description, false)
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

// InitHTTP initialize function
func InitHTTP() {
	authMiddleware, err := CreateAuthMiddleware(os.Getenv("SECRET_KEY"), identifyKey)
	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(-1)
	}

	App = gin.New()
	App.Use(gin.Logger())
	App.Use(gin.Recovery())
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
	}
}
