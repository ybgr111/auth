package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/ybgr111/auth/internal/client/db"
	"github.com/ybgr111/auth/internal/model"
	"github.com/ybgr111/auth/internal/repository"
	"github.com/ybgr111/auth/internal/repository/user/converter"
	userModel "github.com/ybgr111/auth/internal/repository/user/model"
)

const (
	authTable             = "auth"
	idColumn              = "id"
	emailColumn           = "email"
	nameColumn            = "name"
	roleColumn            = "role"
	passwordColumn        = "password"
	passwordConfirmColumn = "password_confirm"
	createdAtColumn       = "created_at"
	updatedAtColumn       = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, user *userModel.User) (int64, error) {
	builderInsert := sq.Insert(authTable).
		PlaceholderFormat(sq.Dollar).
		Columns(emailColumn, nameColumn, roleColumn, passwordColumn, passwordConfirmColumn).
		Values(user.Info.Email, user.Info.Name, user.Info.Role, user.Passwd.Password, user.Passwd.PasswordConfirm).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
		return 0, errors.WithMessage(err, "failed to insert user")
	}

	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.UserPublic, error) {
	builderSelect := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(authTable).
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Where(sq.Eq{"id": id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user userModel.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to scan user")
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) Update(ctx context.Context, req *userModel.User) error {
	builderUpdate := sq.Update(authTable).
		PlaceholderFormat(sq.Dollar).
		Set(nameColumn, req.Info.Name).
		Set(emailColumn, req.Info.Email).
		Set(roleColumn, req.Info.Role).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{"id": req.ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return errors.WithMessage(err, "failed to build query")
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return errors.WithMessage(err, "failed to update user")
	}

	if res.RowsAffected() == 0 {
		return errors.WithMessage(errors.New("failed to update user"), "user not found")
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builderUpdate := sq.Delete(authTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return errors.WithMessage(err, "failed to build query")
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return errors.WithMessage(err, "failed to delete user")
	}

	if res.RowsAffected() == 0 {
		return errors.WithMessage(errors.New("failed to delete user"), "user not found")
	}

	return nil
}
