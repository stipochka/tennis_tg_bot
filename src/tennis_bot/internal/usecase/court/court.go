package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (cu *CourtUsecase) CreateNewCourt(
	ctx context.Context, telegramID int64,
	name, timeOpen, timeClosed, address string,
) (int64, error) {
	log := cu.log.With(slog.String("method", "CreateNewCourt"))
	_, err := cu.repo.CheckIfAdmin(ctx, telegramID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("no admin with such telegramID", slog.Int64("telegramID", telegramID))

			return 0, errNotAnAdmin
		}
		log.Error("failed to check if admin", slog.Any("error", err))
		return 0, err
	}

	log.Debug("admin telegramID", slog.Int64("telegramID", telegramID))

	courtID, err := cu.repo.CreateCourt(ctx, name, timeOpen, timeClosed, address)
	if err != nil {
		log.Error("failed to create court", slog.Any("error", err))

		return 0, err
	}

	log.Debug("court created succsessfully", slog.Int64("courtID", courtID))

	return courtID, nil
}
