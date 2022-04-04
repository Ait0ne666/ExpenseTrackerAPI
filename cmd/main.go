package main

import (
	"expense_tracker/db"
	handlers "expense_tracker/handlers"
	repository "expense_tracker/repositories"
	router "expense_tracker/router"
	"expense_tracker/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	database, err := db.ConnectDatabase()

	if err != nil {
		panic("Error creating database")
	}

	repo := repository.NewRepository(database)

	service := services.NewService(repo)

	handles := handlers.NewHandlers(service)

	router := router.NewRouter(r, handles)

	router.SetupRouter()
	r.Run()

}
