package main

import (
	"log"

	"github.com/SaveliiYam/http_service_with_db/internal/handler"
	"github.com/SaveliiYam/http_service_with_db/internal/middleware"
	"github.com/SaveliiYam/http_service_with_db/internal/repository"
	"github.com/SaveliiYam/http_service_with_db/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	repo := repository.NewUserRepository()
	service := usecase.NewUserService(repo)
	userHandler := handler.NewUserHandler(service)

	router := gin.New()
	router.Use(gin.Recovery(), middleware.RequestLogger())

	userHandler.RegisterRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
