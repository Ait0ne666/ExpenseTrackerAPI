package main

import (
	"expense_tracker/db"
	handlers "expense_tracker/handlers"
	"expense_tracker/pkg/jwt_auth"
	repository "expense_tracker/repositories"
	router "expense_tracker/router"
	"expense_tracker/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// initEnv()
	r := gin.Default()

	database, err := db.ConnectDatabase()

	if err != nil {
		panic("Error creating database")
	}

	repo := repository.NewRepository(database)

	jwt := jwt_auth.NewJwtAuth()
	service := services.NewService(repo, jwt)

	handles := handlers.NewHandlers(service, jwt)

	router := router.NewRouter(r, handles)

	router.SetupRouter()
	r.Run()

}

func initEnv() {
	if err := godotenv.Load("../.env"); err != nil {
		panic("No .env file found")
	}
}
