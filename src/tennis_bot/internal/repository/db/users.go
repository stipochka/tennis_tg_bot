package repository

import (
	"context"
	"fmt"
)

const (
	usersTable    = "users"
	idCol         = "id"
	telegramIDCol = "telegram_id"
	isAdminCol    = "is_admin"
	createdAtCol  = "created_at"
	updatedAtCol  = "updated_at"
)

func (pr *PGRepository) EnsureUser(ctx context.Context, telegramID int64) (int64, error) {
	var userID int64

	query := `INSERT INTO users (telegram_id) VALUES ($1) ON CONFLICT (telegram_id) DO UPDATE SET updated_at = now() RETURNING id`

	err := pr.conn.QueryRow(ctx, query, telegramID).Scan(&userID)

	return userID, err
}

func (pr *PGRepository) MarkAsAdmin(ctx context.Context, telegramID int64) error {
	query := fmt.Sprintf(
		`UPDATE %s SET %s = true WHERE %s = $1;`, usersTable, isAdminCol, telegramIDCol,
	)

	_, err := pr.conn.Exec(ctx, query, telegramID)

	return err
}

func (pr *PGRepository) CheckIfAdmin(ctx context.Context, telegramID int64) (bool, error) {
	var userID int
	err := pr.conn.QueryRow(ctx, queryCheckIfAdmin, telegramID).Scan(&userID)
	if err != nil {
		return false, err
	}

	return true, nil
}
