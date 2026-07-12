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
)

// -- courts - доступные корты
// CREATE TABLE IF NOT EXISTS courts (
//     id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
//     name text NOT NULL,
//     open_time time NOT NULL DEFAULT '07:00',
//     close_time time NOT NULL DEFAULT '23:00',
//     address text NOT NULL,
//     is_active boolean NOT NULL DEFAULT true
// );
