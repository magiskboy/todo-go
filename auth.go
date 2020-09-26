package main

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"time"
)

// LoginInfo json login payload
type LoginInfo struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateAuthMiddleware create auth middleware
func CreateAuthMiddleware(SecretKey, identifyKey string) (*jwt.GinJWTMiddleware, error) {
	// add authenticate middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test",
		Key:         []byte(SecretKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identifyKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identifyKey: v.Email,
					"Name":      v.Name,
					"ID":        v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			var data LoginInfo
			if err := ctx.ShouldBindJSON(&data); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			user, err := GetUserByEmail(data.Email)
			if err != nil || !user.VerifyPassword(data.Password) {
				return "", jwt.ErrFailedAuthentication
			}
			return &user, nil
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		return nil, err
	}
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		return nil, errInit
	}
	return authMiddleware, nil
}
