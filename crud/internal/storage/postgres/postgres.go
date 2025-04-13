package postgres

import (
	"ChatService/crud/internal/domain/models"
	"ChatService/crud/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	// TODO: with postgres
	const op = "storage.postgres.New"

	if _, err := os.Stat("sqlite3://" + storagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: database file does not exist: %w", op, err)
	}

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateMessage(ctx context.Context, uid int64, content string) (int64, error) {
	const op = "storage.postgres.CreateMessage"

	stmt, err := s.db.Prepare("INSERT INTO messages (uid, content) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()

	res, err := stmt.ExecContext(ctx, uid, content)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	mid, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return mid, nil
}

func (s *Storage) GetMessage(ctx context.Context, mid int64) (string, error) {
	const op = "storage.postgres.GetMessage"

	stmt, err := s.db.Prepare("SELECT * FROM messages WHERE id=?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()

	var message models.Message
	if err := stmt.QueryRowContext(ctx, mid).Scan(&message.ID, &message.Content, &message.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrMessageNotExist)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return message.Content, nil
}

func (s *Storage) DeleteMessage(ctx context.Context, mid int64) (bool, error) {
	const op = "storage.postgres.DeleteMessage"
	stmt, err := s.db.Prepare("DELETE from messages WHERE id=?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {

		}
	}()
	res, err := stmt.ExecContext(ctx, mid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrMessageNotExist)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}

func (s *Storage) UpdateMessage(ctx context.Context, mid int64, newContent string) (bool, error) {
	const op = "storage.postgres.UpdateMessage"
	stmt, err := s.db.Prepare("UPDATE messages SET content=? where id=?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()

	res, err := stmt.ExecContext(ctx, newContent, mid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrMessageNotExist)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return n > 0, nil
}
