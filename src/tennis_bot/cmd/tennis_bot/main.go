package main

import (
	"context"
	"fmt"
	repository "tennis_bot/internal/repository/db"
	usecase "tennis_bot/internal/usecase/court"
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

	usecase := usecase.NewCourtUsecase(repo)
	err = usecase.EnsureUser(ctx, 121)
	if err != nil {
		fmt.Println("failed to create user", err)

		return
	}

	repo.MarkAsAdmin(ctx, 121)

	courtID, err := usecase.CreateNewCourt(ctx, 121, "Корт Неймарк №2", "07:00", "23:00", "Большие Овраги 12к18")
	if err != nil {
		fmt.Println("failed to create court", err)
		return
	}

	fmt.Println("Created court id", courtID)

	reservationID, err := usecase.CreateReservation(
		ctx,
		courtID,
		121,
		time.Date(2026, time.July, 5, 14, 0, 0, 0, time.UTC),
		time.Date(2026, time.July, 5, 15, 0, 0, 0, time.UTC),
	)
	if err != nil {
		fmt.Println("error", err)

		return
	}
	fmt.Println(reservationID)
}
