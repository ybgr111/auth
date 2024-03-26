package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ybgr111/auth/internal/service/mocks"
	"github.com/ybgr111/auth/internal/service/user"

	"github.com/ybgr111/platform_common/pkg/db"
	dbMocks "github.com/ybgr111/platform_common/pkg/db/mocks"
	"github.com/ybgr111/platform_common/pkg/db/pg"
	"github.com/ybgr111/platform_common/pkg/db/transaction"

	repoMocks "github.com/ybgr111/auth/internal/repository/mocks"

	converterLogRepo "github.com/ybgr111/auth/internal/repository/log/converter"
)

type deleteUserVariables struct {
	id int64
}

type DeleteUserSuite struct {
	ctx       context.Context
	ctxWithTx context.Context

	suite.Suite

	mc                 *minimock.Controller
	userRepositoryMock *repoMocks.UserRepositoryMock
	logRepositoryMock  *repoMocks.LogRepositoryMock
	fakeTxMock         *mocks.FakeTxMock
	transactorMock     *dbMocks.TransactorMock

	txManagerMock db.TxManager

	deleteUserVariables
}

func TestDeleteUserSuite(t *testing.T) {
	suite.Run(t, new(DeleteUserSuite))
}

func (s *DeleteUserSuite) SetupSuite() {
	s.ctx = context.Background()
	s.mc = minimock.NewController(s.T())

	s.userRepositoryMock = repoMocks.NewUserRepositoryMock(s.mc)
	s.logRepositoryMock = repoMocks.NewLogRepositoryMock(s.mc)
	s.fakeTxMock = mocks.NewFakeTxMock(s.mc)

	s.ctxWithTx = pg.MakeContextTx(s.ctx, s.fakeTxMock)

	s.transactorMock = dbMocks.NewTransactorMock(s.mc)
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	s.transactorMock.BeginTxMock.Expect(s.ctx, txOpts).Return(s.fakeTxMock, nil)

	s.txManagerMock = transaction.NewTransactionManager(s.transactorMock)

	s.id = gofakeit.Int64()
}

func (s *DeleteUserSuite) TestDelete_Success() {
	// Специфичные моки методов.
	s.userRepositoryMock.DeleteMock.Expect(s.ctxWithTx, s.id).Return(nil)
	s.logRepositoryMock.CreateMock.Expect(s.ctxWithTx, converterLogRepo.ToLogCreate("Delete", s.id)).Return(nil)
	s.fakeTxMock.CommitMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	err := service.Delete(s.ctx, s.id)

	// Проверки корректности теста.
	require.Nil(s.T(), nil, err)
}

func (s *DeleteUserSuite) TestDelete_FailDeleteUser() {
	userErr := errors.New("cant delete user")

	s.userRepositoryMock.DeleteMock.Return(userErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	err := service.Delete(s.ctx, s.id)

	require.Error(s.T(), userErr, err)
}

func (s *DeleteUserSuite) TestDelete_FailCreateLog() {
	logErr := errors.New("cant create log")

	s.userRepositoryMock.DeleteMock.Expect(s.ctxWithTx, s.id).Return(nil)
	s.logRepositoryMock.CreateMock.Return(logErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	err := service.Delete(s.ctx, s.id)

	require.Error(s.T(), logErr, err)
}
