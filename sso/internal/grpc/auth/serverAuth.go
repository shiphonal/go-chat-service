package auth

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

type Auth interface {
	Login(ctx context.Context, email, password string, appID int) (string, error)
	Register(ctx context.Context, username, email, password string) (int64, error)
	Logout(ctx context.Context, token string, userID int64) (bool, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	IsModerator(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServiceServer
	auth Auth
}

func RegisterService(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServiceServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	err := validator.LoginValid(req)
	if err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.Email, req.Password, int(req.AppId))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.PermissionDenied, "invalid credentials")
		}
		return nil, status.Error(codes.Unauthenticated, "failed with login")
	}
	return &ssov1.LoginResponse{Token: token}, nil

}

func (s *serverAPI) Logout(ctx context.Context, req *ssov1.LogoutRequest) (*ssov1.LogoutResponse, error) {
	err := validator.LogoutValid(req)
	if err != nil {
		return nil, err
	}

	answer, err := s.auth.Logout(ctx, req.Token, req.UserId)
	return &ssov1.LogoutResponse{Answer: answer}, nil
}

/*func (s *serverAPI) RefreshToken(ctx context.Context, req *ssov1.RefreshTokenRequest) (*ssov1.RefreshTokenResponse, error) {

}*/

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	err := validator.RegisterValid(req)
	if err != nil {
		return nil, err
	}

	id, err := s.auth.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.PermissionDenied, "invalid credentials")
		}
		return nil, status.Error(codes.Unauthenticated, "failed with register")
	}

	return &ssov1.RegisterResponse{UserId: id}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	err := validator.IsAdminValid(req)
	if err != nil {
		return nil, err
	}

	answer, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.PermissionDenied, "invalid credentials")
		}
		return nil, status.Error(codes.Unauthenticated, "failed with isAdmin")
	}
	return &ssov1.IsAdminResponse{IsAdmin: answer}, nil
}

func (s *serverAPI) IsModerator(ctx context.Context, req *ssov1.IsModeratorRequest) (*ssov1.IsModeratorResponse, error) {
	err := validator.IsModValid(req)
	if err != nil {
		return nil, err
	}

	answer, err := s.auth.IsModerator(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.PermissionDenied, "invalid credentials")
		}
		return nil, status.Error(codes.Unauthenticated, "failed with isModerator")
	}
	return &ssov1.IsModeratorResponse{IsMod: answer}, nil
}
