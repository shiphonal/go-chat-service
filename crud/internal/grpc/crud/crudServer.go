package crud

import (
	"ChatService/crud/internal/storage"
	crudv1 "ChatService/protos/gen/go/crud"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CRUD interface {
	GetMessage(ctx context.Context, mid int64) (string, error)
	UpdateMessage(ctx context.Context, mid int64, newContent string) (bool, error)
	SentMessage(ctx context.Context, uid int64, content string) (int64, error)
	DeleteMessage(ctx context.Context, uid int64) (bool, error)
}

type serverCRUD struct {
	crudv1.UnimplementedMessageServer
	crud CRUD
}

func RegisterServer(gRPCServer *grpc.Server, crud CRUD) {
	crudv1.RegisterMessageServer(gRPCServer, &serverCRUD{crud: crud})
}

func (s *serverCRUD) SentMessage(ctx context.Context, req *crudv1.SentMessageRequest) (*crudv1.SentMessageResponse, error) {

	id, err := s.crud.SentMessage(ctx, req.GetUid(), req.GetContent())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to update message")
	}
	return &crudv1.SentMessageResponse{Mid: id}, nil
}

func (s *serverCRUD) DeleteMessage(ctx context.Context, req *crudv1.DeleteMessageRequest) (*crudv1.DeleteMessageResponse, error) {

	answer, err := s.crud.DeleteMessage(ctx, req.GetUid())
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotExist) {
			return nil, status.Error(codes.PermissionDenied, "message not found")
		}
		return nil, status.Error(codes.Unauthenticated, "failed to delete message")
	}
	return &crudv1.DeleteMessageResponse{Status: answer}, nil
}

func (s *serverCRUD) GetMessage(ctx context.Context, req *crudv1.GetMessageRequest) (*crudv1.GetMessageResponse, error) {
	message, err := s.crud.GetMessage(ctx, req.GetUid())
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotExist) {
			return nil, status.Error(codes.PermissionDenied, "message not found")
		}
		return nil, status.Error(codes.Unauthenticated, "failed to get message")
	}
	return &crudv1.GetMessageResponse{Message: message}, nil
}

func (s *serverCRUD) UpdateMessage(ctx context.Context, req *crudv1.UpdateMessageRequest) (*crudv1.UpdateMessageResponse, error) {
	answer, err := s.crud.UpdateMessage(ctx, req.GetUid(), req.GetNewContent())
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotExist) {
			return nil, status.Error(codes.PermissionDenied, "message not found")
		}
		return nil, status.Error(codes.Unauthenticated, "failed to update message")
	}
	return &crudv1.UpdateMessageResponse{Status: answer}, nil
}
