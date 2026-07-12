package repository

import (
	"context"
	"tennis_bot/internal/domain/court"
)

func (pr *PGRepository) CreateCourt(ctx context.Context, name string, openTime, closeTime string, address string) (int64, error) {
	var courtID int64
	err := pr.conn.QueryRow(ctx, queryCreateCourt, name, openTime, closeTime, address).Scan(&courtID)

	return courtID, err
}

func (pr *PGRepository) GetCourts(ctx context.Context) ([]court.Court, error) {
	rows, err := pr.conn.Query(ctx, queryGetCourts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courts []court.Court
	for rows.Next() {
		var (
			id        int64
			name      string
			openTime  string
			closeTime string
			address   string
			isActive  bool
		)

		err := rows.Scan(&id, &name, &openTime, &closeTime, &address, &isActive)
		if err != nil {
			return courts, err
		}

		courts = append(courts, court.Court{
			ID:        id,
			Name:      name,
			OpenHour:  openTime,
			CloseHour: closeTime,
			Address:   address,
			IsActive:  isActive,
		})
	}

	return courts, nil
}

func (pr *PGRepository) GetCourtByID(ctx context.Context, courtID int64) (court.Court, error) {
	var court court.Court
	if err := pr.conn.QueryRow(
		ctx,
		queryGetCourtByID,
		courtID,
	).Scan(
		&court.ID,
		&court.Name,
		&court.OpenHour,
		&court.CloseHour,
		&court.Address,
		&court.IsActive,
	); err != nil {
		return court, err
	}

	return court, nil
}
