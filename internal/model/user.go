package model

import (
	"database/sql"
	"time"
)

type UserPublic struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

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

type UserPassword struct {
	Password        string
	PasswordConfirm string
}
