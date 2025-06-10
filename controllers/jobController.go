package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jobTracker/config"
	"github.com/jobTracker/models"
)

func GetJobs(c *gin.Context) {
	var jobs []models.Job
	config.DB.Find(&jobs)
	c.JSON(http.StatusOK, jobs)
}

func CreateJobs(c *gin.Context) {
	var job models.Job
	if err := c.ShouldBindBodyWithJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&job)
	c.JSON(http.StatusOK, job)
}

func UpdateJobs(c *gin.Context) {
	id := c.Param("id")
	var job models.Job
	if err := config.DB.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	c.BindJSON(&job)
	config.DB.Save(&job)
	c.JSON(http.StatusOK, job)
}

func DeleteJob(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Job{}, id)
	c.Status(http.StatusNoContent)
}
