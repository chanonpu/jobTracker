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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	routes.JobRoutes(r)

	r.Run(":8080")
}
