package main

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

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
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx *gin.Context) interface{} {
			claims := jwt.ExtractClaims(ctx)
			user, _ := GetUserByEmail(claims[identifyKey].(string))
			return &user
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
		Authorizator: func(data interface{}, ctx *gin.Context) bool {
			return true
		},
		Unauthorized: func(ctx *gin.Context, code int, message string) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": message,
			})
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
