package service

import (
	"context"

	"github.com/ybgr111/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, userInfo *model.UserInfo, userPassword *model.UserPassword) (int64, error)
	Get(ctx context.Context, id int64) (*model.UserPublic, error)
	Update(ctx context.Context, id int64, userInfo *model.UserInfo) error
	Delete(ctx context.Context, id int64) error
}
