package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Делаем запрос на вставку записи в таблицу auth
	builderInsert := sq.Insert("auth").
		PlaceholderFormat(sq.Dollar).
		Columns("email", "name", "role", "password", "password_confirm").
		Values(gofakeit.Email(), gofakeit.Name(), gofakeit.Number(0, 2), "qwerty", "qwerty").
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var authID int
	err = pool.QueryRow(ctx, query, args...).Scan(&authID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", authID)

	// Делаем запрос на выборку записей из таблицы auth
	builderSelect := sq.Select("id", "name", "email", "password", "password_confirm", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to select users: %v", err)
	}

	var id int
	var name, email string
	var password, password_confirm string
	var role int
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &name, &email, &password, &password_confirm, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan user: %v", err)
		}

		log.Printf("id: %d, name: %s, email: %s, password: %s, password_confirm: %s, role: %d, created_at: %v, updated_at: %v\n", id, email, name, password, password_confirm, role, createdAt, updatedAt)
	}

	// Делаем запрос на обновление записи в таблице note
	builderUpdate := sq.Update("auth").
		PlaceholderFormat(sq.Dollar).
		Set("name", gofakeit.Name()).
		Set("email", gofakeit.Email()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": authID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	// Делаем запрос на получение измененной записи из таблицы auth
	builderSelectOne := sq.Select("id", "name", "email", "password", "password_confirm", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": authID}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &password, &password_confirm, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select users: %v", err)
	}

	log.Printf("id: %d, name: %s, email: %s, password: %s, password_confirm: %s, role: %d, created_at: %v, updated_at: %v\n", id, email, name, password, password_confirm, role, createdAt, updatedAt)
}
