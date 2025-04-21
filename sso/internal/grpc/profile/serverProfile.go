package profile

import (
	ssov1 "ChatService/protos/gen/go/sso"
	"ChatService/sso/internal/lib/validator"
	"ChatService/sso/internal/services/auth"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Profile interface {
	ChangePassword(ctx context.Context, oldPassword string, password string, id int64) (bool, error)
	ChangeName(ctx context.Context, id int64, newName string) (bool, error)
	ChangeRole(ctx context.Context, password string, idAdmin int64, id int64, newRole int32) (bool, error)
}

type serviceProfile struct {
	ssov1.UnimplementedProfileServer
	profile Profile
}

func RegisterService(grpcServer *grpc.Server, profile Profile) {
	ssov1.RegisterProfileServer(grpcServer, &serviceProfile{profile: profile})
}

func (s *serviceProfile) ChangePassword(ctx context.Context, req *ssov1.ChangePasswordRequest) (*ssov1.ChangePasswordResponse, error) {
	err := validator.ChangePasswordValid(req)
	if err != nil {
		return nil, err
	}

	answer, err := s.profile.ChangePassword(ctx, req.OldPassword, req.NewPassword, req.UserId)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.PermissionDenied, "invalid credentials")
		}
		return nil, status.Error(codes.Unauthenticated, "failed with ChangePassword")
	}
	return &ssov1.ChangePasswordResponse{Success: answer}, nil
}

func (s *serviceProfile) ChangeName(ctx context.Context, req *ssov1.ChangeNameRequest) (*ssov1.ChangeNameResponse, error) {
	err := validator.ChangeNameValid(req)
	if err != nil {
		return nil, err
	}

	answer, err := s.profile.ChangeName(ctx, req.UserId, req.NewName)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.PermissionDenied, "invalid credentials")
		}
		return nil, status.Error(codes.Unauthenticated, "failed with ChangeName")
	}
	return &ssov1.ChangeNameResponse{Success: answer}, nil
}

func (s *serviceProfile) ChangeRole(ctx context.Context, req *ssov1.ChangeRoleRequest) (*ssov1.ChangeRoleResponse, error) {
	err := validator.ChangeRoleValid(req)
	if err != nil {
		return nil, err
	}
	answer, err := s.profile.ChangeRole(ctx, req.Password, req.AdminId, req.UserId, req.NewRole)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.PermissionDenied, "invalid credentials")
		}
		return nil, status.Error(codes.Unauthenticated, "failed with ChangeRole")
	}
	return &ssov1.ChangeRoleResponse{Success: answer}, nil
}
