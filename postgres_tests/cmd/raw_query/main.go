package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4"
)

const (
	dbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Создаем соединение с базой данных
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	// Делаем запрос на вставку записи в таблицу note
	res, err := con.Exec(ctx, "INSERT INTO auth (email, name, role, password, password_confirm) VALUES ($1, $2, $3, $4, $5)", gofakeit.Email(), gofakeit.Name(), gofakeit.Number(0, 2), "qwerty", "qwerty")
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted %d rows", res.RowsAffected())

	// Делаем запрос на выборку записей из таблицы note
	rows, err := con.Query(ctx, "SELECT id, name, email, password, password_confirm, role, created_at, updated_at FROM auth")
	if err != nil {
		log.Fatalf("failed to select user: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		var password, password_confirm string
		var role int
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&id, &name, &email, &password, &password_confirm, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan user: %v", err)
		}

		log.Printf("id: %d, name: %s, email: %s, password: %s, password_confirm: %s, role: %d, created_at: %v, updated_at: %v\n", id, email, name, password, password_confirm, role, createdAt, updatedAt)
	}
}
