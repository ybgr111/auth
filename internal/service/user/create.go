package user

import (
	"context"

	"github.com/pkg/errors"
	"github.com/ybgr111/auth/internal/model"
	converterLogRepo "github.com/ybgr111/auth/internal/repository/log/converter"
	converterUserRepo "github.com/ybgr111/auth/internal/repository/user/converter"
)

func (s *serv) Create(ctx context.Context, userInfo *model.UserInfo, userPassword *model.UserPassword) (int64, error) {
	if userPassword.Password != userPassword.PasswordConfirm {
		return 0, errors.New("passwords dont match")
	}

	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.Create(ctx, converterUserRepo.ToUserCreate(userInfo, userPassword))
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Create(ctx, converterLogRepo.ToLogCreate("Create", id))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
