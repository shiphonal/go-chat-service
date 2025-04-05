package crud

import (
	"ChatService/crud/internal/services/crud"
	"log/slog"
)

func New(log *slog.Logger, cruder crud.MessageCRUDer) *crud.CRUD {
	return &crud.CRUD{
		Log:           log,
		MessageCRUDer: cruder,
	}
}
