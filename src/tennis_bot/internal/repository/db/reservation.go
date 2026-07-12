package repository

import (
	"context"
	"errors"
	"time"

	"tennis_bot/internal/domain/reservation"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	reservationsTable = "reservations"
	courtIDCol        = "court_id"
	userIDCol         = "user_id"
	kindCol           = "kind"
	duringCol         = "during"
	statusCol         = "status"

	reservationNoOverlapCode = "23P01"
)

func (pr *PGRepository) CreateReservation(ctx context.Context, info reservation.Reservation) (int64, error) {
	var reservationID int64
	var pgErr *pgconn.PgError
	err := pr.conn.QueryRow(ctx,
		queryCreateReservation,
		info.CourtID, info.UserID,
		info.Kind,
		info.Start, info.End,
		info.Status,
	).Scan(&reservationID)
	if errors.As(err, &pgErr) {
		if pgErr.Code == reservationNoOverlapCode {
			return reservationID, reservation.ErrReservationOverlap
		}
	}

	return reservationID, err
}

func (pr *PGRepository) ListPending(ctx context.Context, courtID int64) ([]reservation.Reservation, error) {
	var reservations []reservation.Reservation
	rows, err := pr.conn.Query(ctx, queryListPending, courtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id, courtID, userID   int64
			kind, status          string
			start, end, createdAt time.Time
		)

		if err := rows.Scan(
			&id, &courtID,
			&userID, &kind,
			&start, &end,
			&status, &createdAt,
		); err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation.Reservation{
			ID:        id,
			CourtID:   courtID,
			UserID:    userID,
			Kind:      reservation.GetReservationKind(kind),
			Start:     start,
			End:       end,
			Status:    reservation.ReservationStatusPending,
			CreatedAt: createdAt,
		})
	}

	return reservations, nil
}

func (pr *PGRepository) UpdateStatus(ctx context.Context, id int64, status reservation.ReservationStatus) error {
	_, err := pr.conn.Exec(ctx, queryUpdateReservationStatus, status, id)
	return err
}

func (pr *PGRepository) CreateBlockingReservation(
	ctx context.Context,
	courtID, adminID int64,
	start, end time.Time,
) (int64, error) {
	var reservationID int64
	tx, err := pr.conn.Begin(ctx)
	if err != nil {
		return int64(reservationID), err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, queryCancellCrossReservations, courtID, start, end) //TODO: use context with timeout for transaction
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return reservationID, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var cancelledID, userID int64
		err = rows.Scan(&cancelledID, &userID) //TODO: send cancelled reservations to user
	}

	err = tx.QueryRow(ctx, querySetBlockingReservation, courtID, adminID, start, end).Scan(&reservationID)
	if err != nil {
		return reservationID, err
	}

	return reservationID, tx.Commit(ctx)
}
