package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"taskforge/db"
	"taskforge/routers/auth"
)

func main() {
	ctx := context.Background()

	pool, err := db.ConnectDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	authHandler := auth.NewHandler(pool, []byte(jwtKey))

	r := gin.Default()

	authRoutes := r.Group("/auth")
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
