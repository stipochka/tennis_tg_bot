package repository

const (
	queryUpdateReservationStatus = `
		UPDATE reservations
		SET status = $1, reviewed_at = now()
		WHERE id = $2 AND status='pending';
	`
)
