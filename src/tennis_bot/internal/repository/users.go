package repository

import (
	"context"
	"fmt"
	"tennis_bot/internal/domain"
)

const (
	usersTable    = "users"
	idCol         = "id"
	telegramIDCol = "telegram_id"
	isAdminCol    = "is_admin"
	createdAtCol  = "created_at"
	updatedAtCol  = "updated_at"
)

func (pr *PGRepository) EnsureUser(ctx context.Context, telegramID int64) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES ($1) ON CONFLICT (%s) DO UPDATE SET %s = now()
		RETURNING %s, %s, %s, %s;`,
		usersTable, telegramIDCol, telegramIDCol, updatedAtCol,
		idCol, telegramIDCol, isAdminCol, createdAtCol,
	)

	err := pr.conn.QueryRow(ctx, query, telegramID).Scan(&user.ID, &user.TelegramID, &user.IsAdmin, &user.CreatedAt)

	return user, err
}

func (pr *PGRepository) MarkAsAdmin(ctx context.Context, telegramID int64) error {
	query := fmt.Sprintf(
		`UPDATE %s SET %s = true WHERE %s = $1;`, usersTable, isAdminCol, telegramIDCol,
	)

	_, err := pr.conn.Exec(ctx, query, telegramID)

	return err
}
