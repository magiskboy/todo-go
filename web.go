package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

// App is a HTTP application
var App *gin.Engine

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

// LoginRequired authenticate request
func LoginRequired(ctx *gin.Context) (string, error) {
	AuthValue := ctx.Request.Header.Get("Authorization")
	if len(AuthValue) > 1 {
		TokenString := strings.TrimPrefix(AuthValue, "Bearer")
		token, err := jwt.Parse(TokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Method is invalid")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			return "", err
		}
		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// check token expired
			// now := time.Now()
			// if now.After(claims["exp"]) {
			// return "", errors.New("Token is expired")
			// }
			// return claims["email"], nil
			return "", nil
		}
	}
	return "", errors.New("Token is missing")
}

// CheckHealth api check healthy
func CheckHealth(ctx *gin.Context) {
	ctx.Writer.WriteHeader(http.StatusNoContent)
}

// LoginUser login handler
func LoginUser(ctx *gin.Context) {
	var data LoginInfo
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user, err := GetUserByEmail(data.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}
	if user.VerifyPassword(data.Password) {
		token, err := user.GenerateCredential()
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Can't create credential",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		}
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Password is invalid",
		})
	}
}

// CreateUser user login handler
func CreateUser(ctx *gin.Context) {
	var data UserInfo
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	newUser, err := CreateNewUser(data.Email, data.Name, data.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": newUser.ID,
	})
}

// GetTasks get all tasks of a user
func GetTasks(ctx *gin.Context) {

}

// InitHTTP initialize function
func InitHTTP() {
	App = gin.Default()
	App.GET("/", CheckHealth)
	App.POST("/api/auth", LoginUser)
	App.POST("/api/users", CreateUser)
}
