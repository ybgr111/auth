package repository

import (
	"context"

	"github.com/ybgr111/auth/internal/model"
	logModel "github.com/ybgr111/auth/internal/repository/log/model"
	userModel "github.com/ybgr111/auth/internal/repository/user/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *userModel.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.UserPublic, error)
	Update(ctx context.Context, user *userModel.User) error
	Delete(ctx context.Context, id int64) error
}

type LogRepository interface {
	Create(ctx context.Context, log *logModel.Log) error
}
