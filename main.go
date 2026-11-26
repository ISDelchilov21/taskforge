package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"taskforge/db"
	"taskforge/routers"
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
	userHandler := routers.NewUserHandler(pool)
	jobHandler := routers.NewJobHandler(pool)

	r := gin.Default()

	authRoutes := r.Group("/auth")
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)

	api := r.Group("/api", authHandler.AuthRequired())
	api.GET("/users", auth.AdminOnly(), userHandler.GetAllUsers)
	api.GET("/users/:id", userHandler.GetUserByID)
	api.GET("/users/email/:email", userHandler.GetUserbyEmai)
	api.PATCH("/users/:id", userHandler.UpdateUser)
	api.DELETE("/users/:id", userHandler.DleteUser)

	jobs := api.Group("/jobs", auth.AdminOnly())
	jobs.GET("", jobHandler.GetAllJobs)
	jobs.GET("/:id", jobHandler.GetJobByID)
	jobs.POST("", jobHandler.CreateNewJob)
	jobs.PATCH("/:id/status", jobHandler.UpdateJobStatus)
	jobs.DELETE("/:id", jobHandler.DleteJob)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
