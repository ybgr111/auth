package converter

import modelRepo "github.com/ybgr111/auth/internal/repository/log/model"

func ToLogCreate(action string, userId int64) *modelRepo.Log {
	return &modelRepo.Log{
		Action: action,
		UserId: userId,
	}
}
