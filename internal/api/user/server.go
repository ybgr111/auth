package user

import (
	"github.com/ybgr111/auth/internal/service"
	desc "github.com/ybgr111/auth/pkg/note_v1"
)

type Server struct {
	desc.UnimplementedNoteV1Server
	userService service.UserService
}

func NewServer(userService service.UserService) *Server {
	return &Server{
		userService: userService,
	}
}
