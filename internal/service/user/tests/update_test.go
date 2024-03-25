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

	"github.com/ybgr111/auth/internal/client/db"
	"github.com/ybgr111/auth/internal/client/db/pg"
	"github.com/ybgr111/auth/internal/client/db/transaction"
	"github.com/ybgr111/auth/internal/model"
	"github.com/ybgr111/auth/internal/service/mocks"
	"github.com/ybgr111/auth/internal/service/user"

	dbMocks "github.com/ybgr111/auth/internal/client/db/mocks"
	repoMocks "github.com/ybgr111/auth/internal/repository/mocks"

	converterLogRepo "github.com/ybgr111/auth/internal/repository/log/converter"
	converterUserRepo "github.com/ybgr111/auth/internal/repository/user/converter"
)

type updateUserVariables struct {
	id           int64
	name         string
	email        string
	role         int
	password     string
	userInfo     *model.UserInfo
	userPassword *model.UserPassword
}

type UpdateUserSuite struct {
	ctx       context.Context
	ctxWithTx context.Context

	suite.Suite

	mc                 *minimock.Controller
	userRepositoryMock *repoMocks.UserRepositoryMock
	logRepositoryMock  *repoMocks.LogRepositoryMock
	fakeTxMock         *mocks.FakeTxMock
	transactorMock     *dbMocks.TransactorMock

	txManagerMock db.TxManager

	updateUserVariables
}

func TestUpdateUserSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserSuite))
}

func (s *UpdateUserSuite) SetupSuite() {
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
	s.password = gofakeit.Password(true, true, true, false, false, 8)

	s.userInfo = &model.UserInfo{
		Name:  s.name,
		Email: s.email,
		Role:  model.Role(s.role),
	}

	s.userPassword = &model.UserPassword{
		Password:        s.password,
		PasswordConfirm: s.password,
	}
}

func (s *UpdateUserSuite) TestUpdate_Success() {
	// Специфичные моки методов.
	s.userRepositoryMock.UpdateMock.Expect(s.ctxWithTx, converterUserRepo.ToUserUpdate(s.id, s.userInfo)).Return(nil)
	s.logRepositoryMock.CreateMock.Expect(s.ctxWithTx, converterLogRepo.ToLogCreate("Update", s.id)).Return(nil)
	s.fakeTxMock.CommitMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	err := service.Update(s.ctx, s.id, s.userInfo)

	// Проверки корректности теста.
	require.Nil(s.T(), nil, err)
}

func (s *UpdateUserSuite) TestUpdate_FailUpdateUser() {
	userErr := errors.New("cant update user")

	s.userRepositoryMock.UpdateMock.Return(userErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	err := service.Update(s.ctx, s.id, s.userInfo)

	require.Error(s.T(), userErr, err)
}

func (s *UpdateUserSuite) TestUpdate_FailCreateLog() {
	logErr := errors.New("cant create log")

	s.userRepositoryMock.UpdateMock.Expect(s.ctxWithTx, converterUserRepo.ToUserUpdate(s.id, s.userInfo)).Return(nil)
	s.logRepositoryMock.CreateMock.Return(logErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	err := service.Update(s.ctx, s.id, s.userInfo)

	require.Error(s.T(), logErr, err)
}
