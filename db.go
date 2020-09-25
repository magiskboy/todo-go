package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

var db *gorm.DB = nil

// User user model
type User struct {
	gorm.Model
	Name   string `json:"name"`
	Email  string `json:"email" gorm:"unique"`
	Tasks  []Task `json:"tasks" gorm:"foreignKey:ID"`
	PwHash string `gorm:"column:pwhash"`
}

// Task task model
type Task struct {
	gorm.Model
	Name        string `json:"name"`
	Done        bool   `json:"done"`
	Description string `json:"description"`
}

// SetPassword hash raw password and save to db
func (user *User) SetPassword(password string) {
	pwEncrypted := sha256.Sum256([]byte(password))
	user.PwHash = fmt.Sprintf("%x", pwEncrypted)
}

// VerifyPassword check password hash
func (user *User) VerifyPassword(password string) bool {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	return user.PwHash == hash
}

// AddTask append a new task into user
func (user *User) AddTask(name, description string, done bool) (Task, error) {
	NewTask := Task{Name: name, Description: description, Done: done}
	if r := db.Create(&NewTask); r.Error != nil {
		return Task{}, r.Error
	}
	return NewTask, nil
}

// RemoveTask drop a task of the user
func (user *User) RemoveTask(TaskID uint) (Task, error) {
	for _, task := range user.Tasks {
		if task.ID == TaskID {
			db.Delete(&task, TaskID)
			return task, nil
		}
	}
	return Task{}, errors.New("Task not found")
}

// GenerateCredential generate a token
func (user *User) GenerateCredential() (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authenticate"] = true
	atClaims["email"] = user.Email
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return token, err
}

// GetUserByEmail find a user
func GetUserByEmail(email string) (User, error) {
	var user User
	result := db.First(&user, "email = ?", email)
	return user, result.Error
}

// CreateNewUser create a new user
func CreateNewUser(email, name, password string) (User, error) {
	if r := db.Take(&User{}, "email = ?", email); r.RowsAffected > 0 {
		return User{}, errors.New("Email existed")
	}
	NewUser := User{Name: name, Email: email}
	NewUser.SetPassword(password)
	db.Create(&NewUser)
	return NewUser, nil
}

// InitDB create database connection and migrate models
func InitDB() {
	DefaultDatabaseURI := "root:password@tcp(127.0.0.1:3306)/todo?charset=utf8mb4&parseTime=True&loc=Local"
	DatabaseDSN := os.Getenv("DB_DSN")
	if len(DatabaseDSN) == 0 {
		DatabaseDSN = DefaultDatabaseURI
	}

	var err error
	db, err = gorm.Open(mysql.Open(DatabaseDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Task{})
}