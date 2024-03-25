package log

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/ybgr111/auth/internal/client/db"
	"github.com/ybgr111/auth/internal/repository"
	logModel "github.com/ybgr111/auth/internal/repository/log/model"
)

const (
	logTable     = "log_action"
	idColumn     = "id"
	actionColumn = "action"
	userIdColumn = "user_id"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.LogRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, log *logModel.Log) error {
	builderInsert := sq.Insert(logTable).
		PlaceholderFormat(sq.Dollar).
		Columns(actionColumn, userIdColumn).
		Values(log.Action, log.UserId).
		Suffix(fmt.Sprintf("RETURNING %s", idColumn))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "log_repository.Create",
		QueryRaw: query,
	}

	var logID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&logID)
	if err != nil {
		return errors.WithMessage(err, "failed to insert log")
	}

	return nil
}
