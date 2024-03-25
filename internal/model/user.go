package model

import (
	"database/sql"
	"time"
)

// UserPublic публичные данные пользователя.
type UserPublic struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// UserInfo инфо данные пользователя.
type UserInfo struct {
	Name  string
	Email string
	Role  Role
}

type Role int32

const (
	UNSPECIFIED Role = iota
	USER
	ADMIN
)

// UserPassword данные о пароле пользователя.
type UserPassword struct {
	Password        string
	PasswordConfirm string
}
