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

	queryCreateReservation = `
		INSERT INTO reservations (court_id, user_id, kind, during, status)
		VALUES ($1, $2, $3,
			tstzrange(
				$4,
				$5,
				'[)'
			), $6)
		RETURNING id;
	`

	queryListPending = `
		SELECT
			id, court_id, user_id, kind, lower(during), upper(during), status, created_at
	 	FROM reservations WHERE status='pending' AND court_id = $1 ORDER BY created_at;
	`
	queryCreateCourt = `
		INSERT INTO courts (name, open_time, close_time, address)
		VALUES ($1, $2, $3, $4) RETURNING id;
	`

	queryCheckIfAdmin = `
		SELECT id FROM users WHERE telegram_id=$1 and is_admin=true;
	`

	queryGetCourts = `
		SELECT id, name, open_time, close_time, address, is_active
		FROM courts;
	`

	queryGetCourtByID = `
		SELECT id, name, open_time, close_time, address, is_active
		FROM courts WHERE id=$1;
	`

	queryGetAllReservationsByDay = `
		SELECT id, lower(during), upper(during) FROM reservations
	 	WHERE court_id=$1 AND status IN ('pending', 'confirmed') AND during && tstzrange($2, $3, '[)');
	`

	queryApproveReservation = `
		UPDATE reservations
		SET status='confirmed' WHERE id=$1;
	`
)
