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

	err = repo.MartAsAdmin(ctx, 2312)
	if err != nil {
		fmt.Println(err)
	}
	loc, _ := time.LoadLocation("Europe/Moscow")
	reservationID, err := repo.CreateReservation(
		ctx,
		1, 1, domain.ReservationKindBooking,
		time.Date(2026, time.July, 5, 14, 0, 0, 0, loc),
		time.Date(2026, time.July, 5, 15, 0, 0, 0, loc),
		domain.ReservationStatusPending,
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(reservationID)

	reservations, err := repo.ListPending(ctx, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(reservations)
}
