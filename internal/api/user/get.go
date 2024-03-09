package user

import (
	"context"
	"log"

	"github.com/ybgr111/auth/internal/converter"
	desc "github.com/ybgr111/auth/pkg/note_v1"
)

func (i *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userObj, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, role: %v, created_at: %v, updated_at: %v\n", userObj.ID, userObj.Info.Name, userObj.Info.Email, userObj.Info.Role, userObj.CreatedAt, userObj.UpdatedAt)

	return &desc.GetResponse{
		User: converter.ToUserFromService(userObj),
	}, nil
}
