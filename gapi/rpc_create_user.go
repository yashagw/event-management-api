package gapi

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"github.com/yashagw/event-management-api/db/model"
	"github.com/yashagw/event-management-api/pb"
	"github.com/yashagw/event-management-api/util"
	"github.com/yashagw/event-management-api/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(context context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	reqParams := model.CreateUserParams{
		Name:           req.GetName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
	}
	user, err := server.provider.CreateUser(context, reqParams)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "email already exists")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	taskPayload := worker.PayloadSendEmailVerify{
		Email: user.Email,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Timeout(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	err = server.distributor.DistributeTaskSendEmailVerify(context, &taskPayload, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute task: %v", err)
	}

	res := &pb.CreateUserResponse{
		User: &pb.UserResponse{
			Name:              user.Name,
			Email:             user.Email,
			CreatedAt:         timestamppb.New(user.CreatedAt),
			PasswordUpdatedAt: timestamppb.New(user.PasswordUpdatedAt),
			Role:              user.Role.ToProto(),
		},
	}

	return res, nil
}
