package main

import (
	"expense_tracker/db"
	handlers "expense_tracker/handlers"
	"expense_tracker/pkg/jwt_auth"
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

	jwt := jwt_auth.NewJwtAuth()

	token, err := jwt.GenerateJWT("123")

	println(token)

	handles := handlers.NewHandlers(service, jwt)

	router := router.NewRouter(r, handles)

	router.SetupRouter()
	r.Run()

}
