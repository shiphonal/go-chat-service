package crud

import (
	"ChatService/crud/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type CRUD struct {
	log           *slog.Logger
	messageCRUDer MessageCRUDer
}

type MessageCRUDer interface {
	CreateMessage(ctx context.Context, uid int64, content string) (int64, error)
	GetMessage(ctx context.Context, mid int64) (string, error)
	DeleteMessage(ctx context.Context, mid int64) (bool, error)
	UpdateMessage(ctx context.Context, mid int64, newContent string) (bool, error)
}

func (m *CRUD) SentMessage(ctx context.Context, uid int64, content string) (int64, error) {
	const op = "services.crud.SentMessage"
	m.log.With(slog.String("op", op))

	id, err := m.messageCRUDer.CreateMessage(ctx, uid, content)
	if err != nil {
		m.log.Error("Failed to create message")
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (m *CRUD) DeleteMessage(ctx context.Context, uid int64) (bool, error) {
	const op = "services.crud.DeleteMessage"
	m.log.With(slog.String("op", op))
	answer, err := m.messageCRUDer.DeleteMessage(ctx, uid)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotExist) {
			m.log.Error("Message does not exist")
			return false, fmt.Errorf("%s: %w", op, err)
		}
		m.log.Error("Failed to delete message")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return answer, err
}

func (m *CRUD) GetMessage(ctx context.Context, mid int64) (string, error) {
	const op = "services.crud.GetMessage"
	m.log.With(slog.String("op", op))
	content, err := m.messageCRUDer.GetMessage(ctx, mid)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotExist) {
			m.log.Error("Message does not exist")
			return "", fmt.Errorf("%s: %w", op, err)
		}
		m.log.Error("Failed to get message")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return content, err
}

func (m *CRUD) UpdateMessage(ctx context.Context, mid int64, newContent string) (bool, error) {
	const op = "services.crud.UpdateMessage"
	m.log.With(slog.String("op", op))
	answer, err := m.messageCRUDer.UpdateMessage(ctx, mid, newContent)
	if err != nil {
		if errors.Is(err, storage.ErrMessageNotExist) {
			m.log.Error("Message does not exist")
			return false, fmt.Errorf("%s: %w", op, err)
		}
		m.log.Error("Failed to u[date message")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return answer, nil
}
