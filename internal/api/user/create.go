package user

import (
	"context"
	"log"

	"github.com/ybgr111/auth/internal/converter"
	desc "github.com/ybgr111/auth/pkg/note_v1"
)

func (i *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(
		ctx,
		converter.ToUserInfo(req.GetInfo()),
		converter.ToUserPassword(req.GetPasswd()))

	if err != nil {
		return nil, err
	}

	log.Printf("id: %d\n", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
