package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jobTracker/config"
	"github.com/jobTracker/models"
	"github.com/jobTracker/routes"
)

func main() {
	config.Connect()
	config.DB.AutoMigrate(&models.Job{})

	r := gin.Default()

	//enable cors
	r.Use(cors.Default())

	routes.JobRoutes(r)

	r.Run(":8080")
}
