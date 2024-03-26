package user

import (
	"github.com/ybgr111/auth/internal/repository"
	"github.com/ybgr111/auth/internal/service"
	"github.com/ybgr111/platform_common/pkg/db"
)

type serv struct {
	userRepository repository.UserRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

func NewService(
	userRepository repository.UserRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepository: userRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
