package user

import (
	"context"

	"github.com/ybgr111/auth/internal/model"
	converterLogRepo "github.com/ybgr111/auth/internal/repository/log/converter"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.UserPublic, error) {
	var user *model.UserPublic

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		user, errTx = s.userRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Create(ctx, converterLogRepo.ToLogCreate("Get", id))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
