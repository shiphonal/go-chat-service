package crud

import (
	"ChatService/crud/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type CRUD struct {
	Log           *slog.Logger
	MessageCRUDer MessageCRUDer
}

type MessageCRUDer interface {
	CreateMessage(ctx context.Context, uid int64, content string) (int64, error)
	GetMessage(ctx context.Context, mid int64) (string, error)
	DeleteMessage(ctx context.Context, mid int64) (bool, error)
	UpdateMessage(ctx context.Context, mid int64, newContent string) (bool, error)
}

func (m *CRUD) SentMessage(ctx context.Context, uid int64, content string) (int64, error) {
	const op = "services.crud.SentMessage"
	m.Log.With(slog.String("op", op))

	id, err := m.MessageCRUDer.CreateMessage(ctx, uid, content)
	if err != nil {
		m.Log.Error("Failed to create message")
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (m *CRUD) DeleteMessage(ctx context.Context, uid int64) (bool, error) {
	const op = "services.crud.DeleteMessage"
	m.Log.With(slog.String("op", op))
	answer, err := m.MessageCRUDer.DeleteMessage(ctx, uid)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotExist) {
			m.Log.Error("Message does not exist")
			return false, fmt.Errorf("%s: %w", op, err)
		}
		m.Log.Error("Failed to delete message")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return answer, err
}

func (m *CRUD) GetMessage(ctx context.Context, mid int64) (string, error) {
	const op = "services.crud.GetMessage"
	m.Log.With(slog.String("op", op))
	content, err := m.MessageCRUDer.GetMessage(ctx, mid)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotExist) {
			m.Log.Error("Message does not exist")
			return "", fmt.Errorf("%s: %w", op, err)
		}
		m.Log.Error("Failed to get message")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return content, err
}

func (m *CRUD) UpdateMessage(ctx context.Context, mid int64, newContent string) (bool, error) {
	const op = "services.crud.UpdateMessage"
	m.Log.With(slog.String("op", op))
	answer, err := m.MessageCRUDer.UpdateMessage(ctx, mid, newContent)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotExist) {
			m.Log.Error("Message does not exist")
			return false, fmt.Errorf("%s: %w", op, err)
		}
		m.Log.Error("Failed to u[date message")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return answer, nil
}
