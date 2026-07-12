package usecase

import (
	"context"
	"log/slog"
)

func (cu *CourtUsecase) EnsureUser(ctx context.Context, telegramID int64) error {
	log := cu.log.With(slog.String("method", "CreateUser"))

	_, err := cu.repo.EnsureUser(ctx, telegramID)
	if err != nil {
		log.Error("failed method ensure user", slog.Any("error", err))

		return err
	}

	log.Debug("sucsessfully ensured user", slog.Int64("telegramID", telegramID))

	return nil
}

func (cu *CourtUsecase) MarkAsAdmin(ctx context.Context, telegramID int64) error {
	log := cu.log.With(slog.String("method", "MarkAsAdmin"))

	err := cu.repo.MarkAsAdmin(ctx, telegramID)
	if err != nil {
		log.Error("failed to update user", slog.Any("error", err))

		return err
	}

	log.Debug("succsesfully updated user rights", slog.Int64("teleramID", telegramID))

	return nil
}
