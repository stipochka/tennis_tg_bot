package usecase

import (
	"context"
	"log/slog"
	"tennis_bot/internal/domain/reservation"
	"time"

	"os"
)

type Repository interface {
	CreateCourt(ctx context.Context, name string, openTime, closeTime string, address string) (int64, error)
	CreateReservation(ctx context.Context, info reservation.Reservation) (int64, error)
	ListPending(ctx context.Context, courtID int64) ([]reservation.Reservation, error)
	UpdateStatus(ctx context.Context, id int64, status reservation.ReservationStatus) error
	CreateBlockingReservation(
		ctx context.Context,
		courtID, adminID int64,
		start, end time.Time,
	) (int64, error)
	EnsureUser(ctx context.Context, telegramID int64) (int64, error)
	CheckIfAdmin(ctx context.Context, telegramID int64) (bool, error)
	MarkAsAdmin(ctx context.Context, telegramID int64) error
}

type CourtUsecase struct {
	log  slog.Logger
	repo Repository
}

func NewCourtUsecase(repo Repository) CourtUsecase {
	return CourtUsecase{
		log:  *slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		repo: repo,
	}
}
