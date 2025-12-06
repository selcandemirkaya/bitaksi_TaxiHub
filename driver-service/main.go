package main

import (
	_ "bitaksi_taxihub/docs"
	"bitaksi_taxihub/handler"
	"bitaksi_taxihub/repository"
	"bitaksi_taxihub/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	_ = godotenv.Load()

	//Mongo connect
	mongoConn := repository.ConnectDatabase()
	defer mongoConn.CloseMongo()

	//repo
	driverRepo := repository.NewDriverRepositoryMongo(mongoConn.Collection)

	//Service
	driverSvc := service.NewDriverService(driverRepo)

	//Handler
	driverHandler := handler.NewDriverHandler(driverSvc)

	//router
	r := gin.Default()

	//routes
	r.POST("/drivers", driverHandler.CreateDriver)
	r.PUT("/drivers/:id", driverHandler.UpdateDriver)
	r.GET("/drivers", driverHandler.ListDrivers)
	r.GET("/drivers/nearby", driverHandler.NearbyDrivers)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("driver-service running on: 8081")
	r.Run(":8081")
}
