package main

import (
	"context"
	"fmt"
	"tennis_bot/internal/repository"
)

func main() {
	connString := "postgres://court:court@localhost:5432/court"

	repo, err := repository.NewPGRepository(connString)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = repo.EnsureUser(context.Background(), 2312)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = repo.MartAsAdmin(context.Background(), 2312)
	if err != nil {
		fmt.Println(err)
	}
}
