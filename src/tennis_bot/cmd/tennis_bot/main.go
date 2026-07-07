package main

import (
	"context"
	"fmt"
	"tennis_bot/internal/domain"
	"tennis_bot/internal/repository"
	"time"
)

func main() {
	ctx := context.Background()
	connString := "postgres://court:court@localhost:5432/court"

	repo, err := repository.NewPGRepository(connString)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = repo.EnsureUser(ctx, 2312)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = repo.MarkAsAdmin(ctx, 2312)
	if err != nil {
		fmt.Println(err)
	}

	reservationID, err := repo.CreateReservation(
		ctx,
		1, 1, domain.ReservationKindBooking,
		time.Date(2026, time.July, 5, 14, 0, 0, 0, time.UTC),
		time.Date(2026, time.July, 5, 15, 0, 0, 0, time.UTC),
		domain.ReservationStatusPending,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(reservationID)
	reservationID, err = repo.CreateReservation(
		ctx,
		1, 1, domain.ReservationKindBooking,
		time.Date(2026, time.July, 5, 15, 0, 0, 0, time.UTC),
		time.Date(2026, time.July, 5, 17, 0, 0, 0, time.UTC),
		domain.ReservationStatusPending,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Second reservation id", reservationID)

	reservations, err := repo.ListPending(ctx, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(reservations)

	blockingReservationID, err := repo.CreateBlockingReservation(
		ctx,
		1, 1,
		time.Date(2026, time.July, 5, 10, 0, 0, 0, time.UTC),
		time.Date(2026, time.July, 5, 18, 0, 0, 0, time.UTC),
	)
	if err != nil {
		fmt.Println("AA", err)
		return
	}

	fmt.Println(blockingReservationID)
}
