package usecase

import (
	"context"
	"log/slog"
	"tennis_bot/internal/domain/reservation"
	"tennis_bot/internal/domain/validation"
	"time"
)

func (cu *CourtUsecase) CreateReservation(
	ctx context.Context,
	courtID, telegramID int64,
	start, end time.Time,
) (int64, error) {
	log := cu.log.With(slog.String("method", "CreateReservation"))

	err := validation.ValidateTimeBounds(start, end)
	if err != nil {
		log.Info("validation failed", slog.Any("error", err))
		return 0, err
	}
	userID, err := cu.repo.EnsureUser(ctx, telegramID)
	if err != nil {
		log.Error("failed to retrieve userID", slog.Any("error", err))

		return 0, err
	}

	reservationID, err := cu.repo.CreateReservation(ctx, reservation.Reservation{
		CourtID: courtID,
		UserID:  userID,
		Kind:    reservation.ReservationKindBooking,
		Start:   start,
		End:     end,
		Status:  reservation.ReservationStatusPending,
	})
	if err != nil {
		log.Error("failed to create reservation", slog.Any("error", err))

		return 0, err
	}

	log.Debug("reservation succsesfully created", slog.Int64("reservationID", reservationID))

	return reservationID, nil
}
