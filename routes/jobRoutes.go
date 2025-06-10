package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jobTracker/controllers"
)

func JobRoutes(r *gin.Engine) {
	job := r.Group("/jobs")
	{
		job.GET("/", controllers.GetJobs)         // get all job
		job.POST("/", controllers.CreateJobs)     // create
		job.PUT("/:id", controllers.UpdateJobs)   // update
		job.DELETE("/:id", controllers.DeleteJob) // delete
	}
}
