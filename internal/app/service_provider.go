package app

import (
	"context"
	"log"

	"github.com/ybgr111/auth/internal/api/user"
	"github.com/ybgr111/auth/internal/client/db"
	"github.com/ybgr111/auth/internal/client/db/pg"
	"github.com/ybgr111/auth/internal/client/db/transaction"
	"github.com/ybgr111/auth/internal/config"
	"github.com/ybgr111/auth/internal/config/env"
	"github.com/ybgr111/auth/internal/repository"
	logRepository "github.com/ybgr111/auth/internal/repository/log"
	userRepository "github.com/ybgr111/auth/internal/repository/user"
	"github.com/ybgr111/auth/internal/service"
	userService "github.com/ybgr111/auth/internal/service/user"
	"github.com/ybgr111/platform_common/pkg/closer"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository
	logRepository  repository.LogRepository

	userService service.UserService

	server *user.Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: : %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.LogRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) Server(ctx context.Context) *user.Server {
	if s.server == nil {
		s.server = user.NewServer(s.UserService(ctx))
	}

	return s.server
}
