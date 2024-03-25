package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/fatih/color"
	desc "github.com/ybgr111/auth/pkg/note_v1"
)

const (
	address   = "localhost:50052"
	userID    = 12
	userName  = "Nikita Ivanov"
	userEmail = "chertiche@sobaka.budka"
	userPass  = "blabla"
	userRole  = 1
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	create, err := c.Create(ctx, &desc.CreateRequest{
		Info: &desc.UserInfo{
			Name:  userName,
			Email: userEmail,
			Role:  userRole,
		},
		Passwd: &desc.UserPassword{
			Password:        userPass,
			PasswordConfirm: userPass,
		},
	})
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	log.Printf(color.RedString("Сreate user info:\n"), color.GreenString("%+v", create.GetId()))

	get, err := c.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}
	log.Printf(color.RedString("Get user info:\n"), color.GreenString("%+v", get.GetUser()))

	update, err := c.Update(ctx, &desc.UpdateRequest{
		Id: userID,
		Info: &desc.UpdateUserInfo{
			Name:  wrapperspb.String(userName),
			Email: wrapperspb.String(userEmail),
			Role:  userRole,
		},
	})
	if err != nil {
		log.Fatalf("failed to update user by id: %v", err)
	}
	log.Printf(color.RedString("Update user info:\n"), color.GreenString("%+v", update.String()))

	delete, err := c.Delete(ctx, &desc.DeleteRequest{
		Id: userID,
	})
	if err != nil {
		log.Fatalf("failed to delete user by id: %v", err)
	}
	log.Printf(color.RedString("Delete user info:\n"), color.GreenString("%+v", delete.String()))
}
