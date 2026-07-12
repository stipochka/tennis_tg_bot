package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connTimeout = 5 * time.Second
	pingTimeout = 3 * time.Second
)

type PGRepository struct {
	conn *pgxpool.Pool
}

func NewPGRepository(connString string) (*PGRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()
	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), pingTimeout)
	defer pingCancel()
	if err := conn.Ping(pingCtx); err != nil {
		return nil, err
	}

	return &PGRepository{
		conn: conn,
	}, nil
}
