package repository

const (
	queryUpdateReservationStatus = `
		UPDATE reservations
		SET status = $1, reviewed_at = now()
		WHERE id = $2 AND status='pending';
	`
	queryCancellCrossReservations = `
		UPDATE reservations
		SET status='cancelled', cancelled_at=now()
		WHERE court_id = $1 AND status IN ('pending', 'confirmed')
		AND during && tstzrange($2, $3, '[)')
		RETURNING id, user_id;
	`

	 querySetBlockingReservation = `
		INSERT INTO reservations (court_id, user_id, kind, during, status)
		VALUES ($1, $2, 'block', tstzrange($3, $4, '[)'), 'confirmed')
		RETURNING id;
	`
)
