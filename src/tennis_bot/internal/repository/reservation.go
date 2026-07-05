package repository

import (
	"context"
	"fmt"
	"time"

	"tennis_bot/internal/domain"
)

const (
	reservationsTable = "reservations"
	courtIDCol        = "court_id"
	userIDCol         = "user_id"
	kindCol           = "kind"
	duringCol         = "during"
	statusCol         = "status"
)

func (pr *PGRepository) CreateReservation(
	ctx context.Context,
	courtID, userID int64, kind domain.ReservationKind,
	start, end time.Time, status domain.ReservationStatus,
) (int64, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s, %s)
		VALUES ($1, $2, $3,
			tstzrange(
				$4,
				$5,
				'[)'
			), $6)
		RETURNING id;`,
		reservationsTable, courtIDCol, userIDCol,
		kindCol, duringCol, statusCol,
	)

	var reservationID int64
	err := pr.conn.QueryRow(ctx, query, courtID, userID, kind, start, end, status).Scan(&reservationID)

	return reservationID, err
}

func (pr *PGRepository) ListPending(ctx context.Context, courtID int64) ([]domain.Reservation, error) {
	query := fmt.Sprintf(`
		SELECT
			id, court_id, user_id, kind, lower(during), upper(during), status, created_at
	 	FROM %s WHERE status='pending' AND courtID = $1 ORDER BY created_at;`, reservationsTable,
	)

	var reservations []domain.Reservation
	rows, err := pr.conn.Query(ctx, query, courtID)
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

		reservations = append(reservations, domain.Reservation{
			ID:        id,
			CourtID:   courtID,
			UserID:    userID,
			Kind:      domain.GetReservationKind(kind),
			Start:     start,
			End:       end,
			Status:    domain.ReservationStatusPending,
			CreatedAt: createdAt,
		})
	}

	return reservations, nil
}

func (pr *PGRepository) UpdateStatus(ctx context.Context, id int64, status domain.ReservationStatus) error {
	_, err := pr.conn.Exec(ctx, queryUpdateReservationStatus, status, id)
	return err
}

/*
id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
court_id bigint NOT NULL REFERENCES courts(id),
user_id bigint REFERENCES users(id) ON DELETE SET NULL,
kind text NOT NULL DEFAULT 'booking',
during tstzrange NOT NULL,         -- [начало, конец), 14:00-16:00
status text NOT NULL DEFAULT 'pending',
created_at timestamptz NOT NULL DEFAULT now(),
reviwed_at timestamptz,
cancelled_at timestamptz,
*/
