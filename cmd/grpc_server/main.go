package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ybgr111/auth/internal/config"
	"github.com/ybgr111/auth/internal/config/env"
	desc "github.com/ybgr111/auth/pkg/note_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedNoteV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	builderInsert := sq.Insert("auth").
		PlaceholderFormat(sq.Dollar).
		Columns("email", "name", "role", "password", "password_confirm").
		Values(req.Info.Email, req.Info.Name, req.Info.Role, req.Passwd.Password, req.Passwd.PasswordConfirm).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to build query")
	}

	var authID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&authID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to insert user")
	}

	log.Printf("User info: %v", req.GetInfo())
	log.Printf("User passwd: %v", "ne napishy, ver' mne")

	return &desc.CreateResponse{
		Id: authID,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var id int64
	var name, email string
	var role desc.RoleType
	var createdAt time.Time
	var updatedAt sql.NullTime

	row := s.pool.QueryRow(ctx, query, args...)
	err = row.Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to scan user")
	}

	log.Printf("User id: %d", req.GetId())

	return &desc.GetResponse{
		User: &desc.User{
			Id: id,
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  role,
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt.Time),
		},
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	builderUpdate := sq.Update("auth").
		PlaceholderFormat(sq.Dollar).
		Set("name", req.Info.Name.Value).
		Set("email", req.Info.Email.Value).
		Set("role", req.Info.Role).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to build query")
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update user")
	}

	if res.RowsAffected() == 0 {
		log.Printf("sdsd")
		return nil, errors.WithMessage(errors.New("failed to update user"), "user not found")
	}

	log.Printf("User id: %v", req.GetId())
	log.Printf("User info: %v", req.GetInfo())

	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	builderUpdate := sq.Delete("auth").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to build query")
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to delete user")
	}

	if res.RowsAffected() == 0 {
		log.Printf("sdsd")
		return nil, errors.WithMessage(errors.New("failed to delete user"), "user not found")
	}

	log.Printf("User id: %v", req.GetId())

	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: #{err}")
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: #{err}")
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
