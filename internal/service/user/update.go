package user

import (
	"context"

	"github.com/ybgr111/auth/internal/model"
	converterLogRepo "github.com/ybgr111/auth/internal/repository/log/converter"
	converterUserRepo "github.com/ybgr111/auth/internal/repository/user/converter"
)

func (s *serv) Update(ctx context.Context, id int64, userInfo *model.UserInfo) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Update(ctx, converterUserRepo.ToUserUpdate(id, userInfo))
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Create(ctx, converterLogRepo.ToLogCreate("Update", id))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
