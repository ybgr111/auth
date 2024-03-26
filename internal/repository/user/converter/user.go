package converter

import (
	"github.com/ybgr111/auth/internal/model"
	modelRepo "github.com/ybgr111/auth/internal/repository/user/model"
)

// ToUserInfoFromRepo конвертер репо - сервис по возврату информации о пользователе.
func ToUserFromRepo(user *modelRepo.User) *model.UserPublic {
	return &model.UserPublic{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(&user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserInfoFromRepo конвертер репо - сервис по возврату личной информации о пользователе.

func ToUserInfoFromRepo(info *modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  model.Role(info.Role),
	}
}

// ToUserUpdate конвертер сервис - репо по созданию пользователя.
func ToUserCreate(
	userInfo *model.UserInfo,
	userPassword *model.UserPassword,
) *modelRepo.User {
	return &modelRepo.User{
		Info: modelRepo.UserInfo{
			Name:  userInfo.Name,
			Email: userInfo.Email,
			Role:  modelRepo.Role(userInfo.Role),
		},
		Passwd: modelRepo.UserPassword(*userPassword),
	}
}

// ToUserUpdate конвертер сервис - репо по обновлению пользователя.
func ToUserUpdate(
	id int64,
	userInfo *model.UserInfo,
) *modelRepo.User {
	return &modelRepo.User{
		ID: id,
		Info: modelRepo.UserInfo{
			Name:  userInfo.Name,
			Email: userInfo.Email,
			Role:  modelRepo.Role(userInfo.Role),
		},
	}
}
