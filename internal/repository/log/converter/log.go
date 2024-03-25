package converter

import modelRepo "github.com/ybgr111/auth/internal/repository/log/model"

// ToLogCreate конвертер данных сервис - репо по созданию лога.
func ToLogCreate(action string, userId int64) *modelRepo.Log {
	return &modelRepo.Log{
		Action: action,
		UserId: userId,
	}
}
