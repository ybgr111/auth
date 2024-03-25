package model

type Log struct {
	ID     int64  `db:"id"`
	Action string `db:"action"`
	UserId int64  `db:"user_id"`
}
