package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type TodoTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (testSuite *TodoTestSuite) SetupTest() {
	os.Setenv("ENV", "testing")
	gin.SetMode("release")

	InitDB()
	testSuite.router = SetupRouter()
}

func (testSuite *TodoTestSuite) TestCheckHealth() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	testSuite.router.ServeHTTP(w, req)

	assert.Equal(testSuite.T(), http.StatusNoContent, w.Code)
}

func (testSuite *TodoTestSuite) TestUnauthorizeRequest() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/private/tasks", nil)
	testSuite.router.ServeHTTP(w, req)

	assert.Equal(testSuite.T(), http.StatusUnauthorized, w.Code)
}

func (testSuite *TodoTestSuite) TestRegisterUser() {
	w := httptest.NewRecorder()
	userData := CreateUserRequest{
		Name:     "lmao",
		Email:    "user@email.com",
		Password: "123456",
	}
	body, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/public/users", bytes.NewReader(body))
	testSuite.router.ServeHTTP(w, req)

	assert.Equal(testSuite.T(), http.StatusOK, w.Code)
	var resp CreateUserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	var user User
	result := db.First(&user, resp.UserID)
	assert.Nil(testSuite.T(), result.Error)
	assert.Equal(testSuite.T(), userData.Email, user.Email)
	assert.Equal(testSuite.T(), userData.Name, user.Name)
	assert.Nil(testSuite.T(), err)
}

func (testSuite *TodoTestSuite) TestLoginUser() {
	password := "123456"
	user := User{Name: "Lmao", Email: "user@email.com"}
	user.SetPassword(password)
	db.Create(&user)
	loginInfo := LoginInfo{Email: user.Email, Password: password}
	loginData, _ := json.Marshal(loginInfo)
	req, _ := http.NewRequest("POST", "/public/login", bytes.NewReader(loginData))
	w := httptest.NewRecorder()
	testSuite.router.ServeHTTP(w, req)

	assert.Equal(testSuite.T(), http.StatusOK, w.Code)
}

func (testSuite *TodoTestSuite) TestGetTasks() {
}

func TestTodoTestSuite(t *testing.T) {
	suite.Run(t, new(TodoTestSuite))
}
