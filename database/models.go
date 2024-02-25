package database

import (
	"database/sql"
	"time"
)

type App struct {
	ID    string
	AppID string
}

type User struct {
	ID             string
	Name           string
	Email          string
	ObjectID       string
	TID            string
	Username       string
	ProfilePicture sql.NullString
	Admin          int
	Disabled       int
}

type UserFeedback struct {
	ID        string
	UserID    string
	Category  string
	Title     string
	Feedback  string
	Status    string
	CreatedAt string
	UpdatedAt string
}

type UserActivity struct {
	ID        string
	UserID    string
	Data      string
	CreatedAt time.Time
}

type UserStatusSettings struct {
	ID                string
	DefaultUserStatus string
}
