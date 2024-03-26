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

	"github.com/ybgr111/auth/internal/model"
	"github.com/ybgr111/auth/internal/service/mocks"
	"github.com/ybgr111/auth/internal/service/user"

	"github.com/ybgr111/platform_common/pkg/db"
	dbMocks "github.com/ybgr111/platform_common/pkg/db/mocks"
	"github.com/ybgr111/platform_common/pkg/db/pg"
	"github.com/ybgr111/platform_common/pkg/db/transaction"

	repoMocks "github.com/ybgr111/auth/internal/repository/mocks"

	converterLogRepo "github.com/ybgr111/auth/internal/repository/log/converter"
	converterUserRepo "github.com/ybgr111/auth/internal/repository/user/converter"
)

type createUserVariables struct {
	id           int64
	name         string
	email        string
	role         int
	password     string
	userInfo     *model.UserInfo
	userPassword *model.UserPassword
}

type CreateUserSuite struct {
	ctx       context.Context
	ctxWithTx context.Context

	suite.Suite

	mc                 *minimock.Controller
	userRepositoryMock *repoMocks.UserRepositoryMock
	logRepositoryMock  *repoMocks.LogRepositoryMock
	fakeTxMock         *mocks.FakeTxMock
	transactorMock     *dbMocks.TransactorMock

	txManagerMock db.TxManager

	createUserVariables
}

func TestCreateUserSuite(t *testing.T) {
	suite.Run(t, new(CreateUserSuite))
}

func (s *CreateUserSuite) SetupSuite() {
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

func (s *CreateUserSuite) TestCreate_Success() {
	// Специфичные моки методов.
	s.userRepositoryMock.CreateMock.Expect(s.ctxWithTx, converterUserRepo.ToUserCreate(s.userInfo, s.userPassword)).Return(s.id, nil)
	s.logRepositoryMock.CreateMock.Expect(s.ctxWithTx, converterLogRepo.ToLogCreate("Create", s.id)).Return(nil)
	s.fakeTxMock.CommitMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	newID, err := service.Create(s.ctx, s.userInfo, s.userPassword)

	// Проверки корректности теста.
	require.Equal(s.T(), nil, err)
	require.Equal(s.T(), s.id, newID)
}

func (s *CreateUserSuite) TestCreate_FailCreateUser() {
	userErr := errors.New("cant create user")

	s.userRepositoryMock.CreateMock.Return(0, userErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	_, err := service.Create(s.ctx, s.userInfo, s.userPassword)

	require.Error(s.T(), userErr, err)
}

func (s *CreateUserSuite) TestCreate_FailCreateLog() {
	logErr := errors.New("cant create log")

	s.userRepositoryMock.CreateMock.Expect(s.ctxWithTx, converterUserRepo.ToUserCreate(s.userInfo, s.userPassword)).Return(0, nil)
	s.logRepositoryMock.CreateMock.Return(logErr)
	s.fakeTxMock.RollbackMock.Return(nil)

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	_, err := service.Create(s.ctx, s.userInfo, s.userPassword)

	require.Error(s.T(), logErr, err)
}

func (s *CreateUserSuite) TestCreate_FailEqualPasswords() {
	passwdErr := errors.New("passwords dont match")

	badUserPassword := &model.UserPassword{
		Password:        "12345678",
		PasswordConfirm: "87654321",
	}

	service := user.NewService(s.userRepositoryMock, s.logRepositoryMock, s.txManagerMock)

	_, err := service.Create(s.ctx, s.userInfo, badUserPassword)

	require.Error(s.T(), passwdErr, err)
}
