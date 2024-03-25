package user

import (
	"context"

	converterLogRepo "github.com/ybgr111/auth/internal/repository/log/converter"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Create(ctx, converterLogRepo.ToLogCreate("Delete", id))
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
