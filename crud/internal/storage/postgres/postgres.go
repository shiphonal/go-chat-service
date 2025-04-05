package postgres

import (
	"ChatService/crud/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	const op = "storage.postgres.New"
	password := os.Getenv("POSTGRES_PASSWORD")

	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

const (
	host   = "localhost"
	port   = "5432"
	user   = "postgres"
	dbname = "postgres"
)

func (s *Storage) CreateMessage(ctx context.Context, uid int64, content string) (int64, error) {
	const op = "storage.postgres.CreateMessage"
	// TODO: smth with postgres

	stmt, err := s.db.Prepare("INSERT INTO postgres (uid, content) VALUES (?, ?)")
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

	stmt, err := s.db.Prepare("SELECT content FROM messages WHERE id=?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()

	res, err := stmt.QueryContext(ctx, mid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrMessageNotExist)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var message string
	if err := res.Scan(&message); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return message, nil
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
