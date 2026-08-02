package usecase

import (
	"context"
	"log/slog"
	"tennis_bot/internal/domain/reservation"
	"tennis_bot/internal/domain/validation"
	"time"
)

const (
	openHour  = 7
	closeHour = 23 //TODO: replace with query to court
)

func (cu *CourtUsecase) CreateReservation(
	ctx context.Context,
	courtID, telegramID int64,
	start, end time.Time,
) (int64, error) {
	log := logWithSource(cu.log, "CreateReservation")

	err := validation.ValidateTimeBounds(time.Now(), start, end)
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

func (cu *CourtUsecase) ListAvaliableSlotsByDay(ctx context.Context, hours int, date time.Time) ([]reservation.Slot, error) {
	log := logWithSource(cu.log, "ListAvaliableSlotsByDay")

	searchedDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	slots, err := cu.repo.ListAllDaySlots(ctx, 1, searchedDay)
	if err != nil {
		log.Error("failed to list all day slots", slog.Any("error", err))

		return nil, err
	}

	return cu.calcAvailableDaySlots(hours, slots), nil
}

func (cu *CourtUsecase) GetAllReservationsByDay(ctx context.Context, date time.Time) ([]reservation.DaySlot, error) {
	log := logWithSource(cu.log, "GetAllReservationsByDay")

	searchedDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	slots, err := cu.repo.ListAllDaySlots(ctx, 1, searchedDay)
	if err != nil {
		log.Error("failed to list all day slots", slog.Any("error", err))

		return nil, err
	}

	return slots, nil
}

func (cu *CourtUsecase) ApproveReservation(ctx context.Context, id int) error {
	return cu.repo.ApproveReservation(ctx, id)
}

func (cu *CourtUsecase) calcAvailableDaySlots(hours int, takenSlots []reservation.DaySlot) []reservation.Slot {
	var allSlots []reservation.Slot
	start := openHour
	for start+hours <= closeHour {
		isTaken := false
		slotCandidate := reservation.Slot{
			Start: start,
			End:   start + hours,
		}

		for _, slot := range takenSlots {
			start := slot.Start.Hour()
			end := slot.End.Hour()

			if max(slotCandidate.Start, start) < min(slotCandidate.End, end) {
				isTaken = true
				break
			}
		}

		if !isTaken {
			allSlots = append(allSlots, slotCandidate)
		}

		start += 1
	}

	return allSlots
}
