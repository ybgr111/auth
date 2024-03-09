package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ybgr111/auth/internal/model"
	desc "github.com/ybgr111/auth/pkg/note_v1"
)

func ToUserFromService(user *model.UserPublic) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Info:      ToUserInfoFromService(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserInfoFromService(info model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  desc.RoleType(info.Role),
	}
}

func ToUserInfo(userInfo *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  userInfo.Name,
		Email: userInfo.Email,
		Role:  model.Role(userInfo.Role),
	}
}

func ToUserPassword(userPassword *desc.UserPassword) *model.UserPassword {
	return &model.UserPassword{
		Password:        userPassword.Password,
		PasswordConfirm: userPassword.PasswordConfirm,
	}
}

func ToUpdateUserInfo(userInfo *desc.UpdateUserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  userInfo.Name.Value,
		Email: userInfo.Email.Value,
		Role:  model.Role(userInfo.Role),
	}
}
