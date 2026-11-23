package main

import (
	"context"
	"fmt"
	"log"

	"taskforge/db"
)

func main() {
	ctx := context.Background()

	pool, err := db.ConnectDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	fmt.Println("Database connected")
}
