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
		`INSERT INTO %s (%s, %s, %s, %s)
		VALUES ($1, $2, $3, tsrzrange($4, $5, '[)')
		RETURNING id;`,
		reservationsTable, courtIDCol, userIDCol,
		kindCol, duringCol,
	)

	var reservationID int64
	err := pr.conn.QueryRow(ctx, query, courtID, userID, kind, start, end).Scan(&reservationID)

	return reservationID, err
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
