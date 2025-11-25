package server

import (
	"context"
	"time"

	"github.com/Kars07/rest-grpc-demo/models"
	pb "github.com/Kars07/rest-grpc-demo/proto"
	"github.com/Kars07/rest-grpc-demo/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	service *service.UserService
}

func NewUserServer(service *service.UserService) *UserServer {
	return &UserServer{service: service}
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.service.GetUserByID(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return userToProto(user), nil
}

func (s *UserServer) GetAllUsers(ctx context.Context, req *pb.Empty) (*pb.UserListResponse, error) {
	users, err := s.service.GetAllUsers()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoUsers []*pb.UserResponse
	for _, user := range users {
		protoUsers = append(protoUsers, userToProto(&user))
	}

	return &pb.UserListResponse{Users: protoUsers}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	createReq := &models.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	user, err := s.service.CreateUser(createReq)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return userToProto(user), nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	updateReq := &models.UpdateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	user, err := s.service.UpdateUser(req.Id, updateReq)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return userToProto(user), nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
	err := s.service.DeleteUser(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.Empty{}, nil
}

func (s *UserServer) StreamUsers(req *pb.Empty, stream pb.UserService_StreamUsersServer) error {
	users, err := s.service.GetAllUsers()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, user := range users {
		if err := stream.Send(userToProto(&user)); err != nil {
			return err
		}
		// Simulate streaming delay
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// Helper function to convert User model to proto
func userToProto(user *models.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
