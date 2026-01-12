package main

import (
	"log"
	"my-sql/internal/config"
	"my-sql/internal/db"
	"my-sql/internal/handler"
	"my-sql/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	mysqlDB, err := db.NewMySQLConnection(
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(mysqlDB)
	userHandler := handler.NewUserHandler(userRepo)

	r := gin.Default()
	r.GET("/users", userHandler.GetUsers)

	r.Run(":8080")
}
