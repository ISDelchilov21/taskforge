package routers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"taskforge/db"
	"taskforge/models"
)

type JobHandler struct {
	pool *pgxpool.Pool
}

func NewJobHandler(pool *pgxpool.Pool) *JobHandler {
	return &JobHandler{pool: pool}
}

func (h *JobHandler) GetAllJobs(c *gin.Context) {
	jobs, err := db.ListJobs(c.Request.Context(), h.pool)
	if err != nil {
		log.Printf("list jobs: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) CreateNewJob(c *gin.Context) {
	var req models.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job := models.Job{
		Type:           strings.TrimSpace(req.Type),
		Description:    strings.TrimSpace(req.Description),
		Payload:        req.Payload,
		PriorityWeight: req.PriorityWeight,
	}

	if err := db.CreateJob(c.Request.Context(), h.pool, &job); err != nil {
		log.Printf("create job: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, job)
}

func (h *JobHandler) GetJobByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	job, err := db.GetJob(c.Request.Context(), h.pool, id)
	if errors.Is(err, db.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}
	if err != nil {
		log.Printf("get job %d: %v", id, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) UpdateJobStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	var req models.UpdateJobStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "status must be one of: pending, running, done, failed"})
		return
	}

	err = db.UpdateJobStatus(c.Request.Context(), h.pool, id, req.Status)
	if errors.Is(err, db.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}
	if err != nil {
		log.Printf("update job %d: %v", id, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	job, err := db.GetJob(c.Request.Context(), h.pool, id)
	if errors.Is(err, db.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}
	if err != nil {
		log.Printf("get job %d: %v", id, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) DleteJob(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	err = db.DeleteJob(c.Request.Context(), h.pool, id)
	if errors.Is(err, db.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}
	if err != nil {
		log.Printf("delete job %d: %v", id, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Status(http.StatusNoContent)
}
