package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"github.com/ybgr111/auth/internal/api/user"
	"github.com/ybgr111/auth/internal/model"
	"github.com/ybgr111/auth/internal/service"
	serviceMocks "github.com/ybgr111/auth/internal/service/mocks"
	desc "github.com/ybgr111/auth/pkg/note_v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUpdate(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id         = gofakeit.Int64()
		name       = gofakeit.FirstName()
		email      = gofakeit.Email()
		role       = gofakeit.Number(0, 2)
		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateRequest{
			Id: id,
			Info: &desc.UpdateUserInfo{
				Name:  wrapperspb.String(name),
				Email: wrapperspb.String(email),
				Role:  desc.RoleType(role),
			},
		}

		userInfo = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  model.Role(role),
		}

		res = &empty.Empty{}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *empty.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, id, userInfo).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, id, userInfo).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := user.NewServer(userServiceMock)

			res, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
