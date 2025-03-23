package auth

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
	ssov1 "protos/gen/go/sso"
)

type Auth interface {
	Login(ctx context.Context, email, password string, appID int) (string, error)
	Register(ctx context.Context, username, email, password string) (int64, error)
	Logout(ctx context.Context)
}

func RegisterService(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func (s *serverAPI) Login(ctx context.Context) {
	log := slog.Logger{}
	log.Info("starting gRPC auth")

}

func (s *serverAPI) Logout() {}

func (s *serverAPI) RefreshToken() {}

func (s *serverAPI) Register() {}
