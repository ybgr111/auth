package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64        `db:"id"`
	Info      UserInfo     `db:""`
	Passwd    UserPassword `db:""`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type UserInfo struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  Role   `db:"role"`
}

type Role int32

const (
	UNSPECIFIED Role = iota
	USER
	ADMIN
)

type UserPassword struct {
	Password        string `db:"password"`
	PasswordConfirm string `db:"password_confirm"`
}
