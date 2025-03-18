package sqlite

import (
	"ChatService/sso/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "sqlite.New"
	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, username, email, password string) (int64, error) {
	const op = "sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()
	res, err := stmt.ExecContext(ctx, username, email, password)
	if err != nil {
		// TODO: retry errors
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (models.User, error) {
	const op = "sqlite.GetUser"

	stmt, err := s.db.Prepare("SELECT * FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()
	var user models.User
	if err := stmt.QueryRowContext(ctx, email).
		Scan(&user.ID, &user.UserName, &user.Email, &user.PassHash); err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) SaveApp(ctx context.Context, name, secret string) (int64, error) {
	const op = "sqlite.SaveApp"

	stmt, err := s.db.Prepare("INSERT INTO app (name, secret) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()
	res, err := stmt.ExecContext(ctx, name, secret)
	if err != nil {
		// TODO: retry errors
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetApp(ctx context.Context, id int) (models.App, error) {
	const op = "sqlite.GetApp"

	stmt, err := s.db.Prepare("SELECT * FROM app WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()
	var app models.App
	if err := stmt.QueryRowContext(ctx, id).
		Scan(&app.ID, &app.AppName, &app.Secret); err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, id int64) (bool, error) {
	const op = "sqlite.IsAdmin"
	stmt, err := s.db.Prepare("SELECT role FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()

	var role string
	if err := stmt.QueryRowContext(ctx, id).Scan(&role); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return role == "admin", nil
}

func (s *Storage) IsModerator(ctx context.Context, id int64) (bool, error) {
	const op = "sqlite.IsAdmin"
	stmt, err := s.db.Prepare("SELECT role FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()

	var role string
	if err := stmt.QueryRowContext(ctx, id).Scan(&role); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return role == "mod", nil
}

func (s *Storage) UpdatePassword(ctx context.Context, oldPassword, newPassword string) (bool, error) {
	const op = "sqlite.UpdatePassword"
	return false, nil
}

func (s *Storage) UpdateName(ctx context.Context, id int64, newName string) (bool, error) {
	const op = "sqlite.UpdateName"
	stmt, err := s.db.Prepare("UPDATE users SET name = ? WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()
	if _, err := stmt.ExecContext(ctx, newName, id); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}

func (s *Storage) ChangeRole(ctx context.Context, idPerson, id int64, newRole string) (bool, error) {
	const op = "sqlite.ChangeRole"
	stmt, err := s.db.Prepare("UPDATE users SET role = ? WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()
	if success, err := s.IsAdmin(ctx, idPerson); success && err != nil {
		if _, err := stmt.ExecContext(ctx, newRole, id); err != nil {
			return false, fmt.Errorf("%s: %w", op, err)
		}
		return true, nil
	} else {
		return false, nil
	}
}
