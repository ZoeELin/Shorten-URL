package main

import (
	"github.com/gin-gonic/gin"
	"golangAPI/api_service"
	"golangAPI/db_service"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	db_service.DatabaseConnect()
	db_service.CreateTable()
	router := gin.Default()

	router.POST("/api/v1/urls", api_service.PostUrl)
	router.GET("/:url_id", api_service.GetUrl)

	router.Run(":8000")

	db_service.CloseDatabase()
}