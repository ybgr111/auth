package user

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ybgr111/auth/internal/converter"
	desc "github.com/ybgr111/auth/pkg/note_v1"
)

func (i *Server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	err := i.userService.Update(
		ctx,
		req.GetId(),
		converter.ToUpdateUserInfo(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	log.Printf("user with id: %d updated\n", req.GetId())

	return &empty.Empty{}, nil
}
