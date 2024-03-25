package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ybgr111/auth/internal/client/db"
	"github.com/ybgr111/auth/internal/client/db/pg"
	"github.com/ybgr111/auth/internal/client/db/transaction"
	"github.com/ybgr111/auth/internal/model"
	"github.com/ybgr111/auth/internal/service/mocks"
	"github.com/ybgr111/auth/internal/service/user"

	dbMocks "github.com/ybgr111/auth/internal/client/db/mocks"
	repoMocks "github.com/ybgr111/auth/internal/repository/mocks"

	converterLogRepo "github.com/ybgr111/auth/internal/repository/log/converter"
)

type getUserVariables struct {
	id         int64
	name       string
	email      string
	role       int
	createdAt  time.Time
	updatedAt  time.Time
	userInfo   *model.UserInfo
	userPublic *model.UserPublic
}

type GetUserSuite struct {
	ctx       context.Context
	ctxWithTx context.Context

	suite.Suite

	mc                 *minimock.Controller
	userRepositoryMock *repoMocks.UserRepositoryMock
	logRepositoryMock  *repoMocks.LogRepositoryMock
	fakeTxMock         *mocks.FakeTxMock
	transactorMock     *dbMocks.TransactorMock

	txManagerMock db.TxManager

	getUserVariables
}

func TestGetUserSuite(t *testing.T) {
	suite.Run(t, new(GetUserSuite))
}

func (s *GetUserSuite) SetupSuite() {
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
	s.name = gofakeit.FirstName()
	s.email = gofakeit.Email()
	s.role = gofakeit.Number(0, 2)
	s.createdAt = gofakeit.Date()
	s.updatedAt = gofakeit.Date()

	s.userInfo = &model.UserInfo{
		Name:  s.name,
		Email: s.email,
		Role:  model.Role(s.role),
	}

	s.userPublic = &model.UserPublic{
		ID: s.id,
		Info: model.UserInfo{
			Name:  s.name,
			Email: s.email,
			Role:  model.Role(s.role),
		},
		CreatedAt: s.createdAt,
		UpdatedAt: sql.NullTime{Time: s.updatedAt},
	}
}

func (s *GetUserSuite) TestGet_Success() {
	// Специфичные моки методов.
	s.userRepositoryMock.GetMock.Expect(s.ctxWithTx, s.id).Return(s.userPublic, nil)
	s.logRepositoryMock.CreateMock.Expect(s.ctxWithTx, converterLogRepo.ToLogCreate("Get", s.id)).Return(nil)
	s.fakeTxMock.CommitMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	data, err := service.Get(s.ctx, s.id)

	// Проверки корректности теста.
	require.Nil(s.T(), nil, err)
	require.Equal(s.T(), s.userPublic, data)
}

func (s *GetUserSuite) TestGet_FailGetUser() {
	userErr := errors.New("cant get user")

	s.userRepositoryMock.GetMock.Return(nil, userErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	_, err := service.Get(s.ctx, s.id)

	require.Error(s.T(), userErr, err)
}

func (s *GetUserSuite) TestDelete_FailCreateLog() {
	logErr := errors.New("cant create log")

	s.userRepositoryMock.GetMock.Expect(s.ctxWithTx, s.id).Return(s.userPublic, nil)
	s.logRepositoryMock.CreateMock.Return(logErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	_, err := service.Get(s.ctx, s.id)

	require.Error(s.T(), logErr, err)
}
