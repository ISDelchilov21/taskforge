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

type UserHandler struct {
	pool *pgxpool.Pool
}

func NewUserHandler(pool *pgxpool.Pool) *UserHandler {
	return &UserHandler{pool: pool}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := db.ListUsers(c.Request.Context(), h.pool)
	if err != nil {
		log.Printf("list users: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	if c.GetInt("userID") != id && c.GetString("role") != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you can only view your own account"})
		return
	}

	user, err := db.GetUserByID(c.Request.Context(), h.pool, id)
	if errors.Is(err, db.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		log.Printf("get user %d: %v", id, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserbyEmai(c *gin.Context) {
	email := strings.ToLower(strings.TrimSpace(c.Param("email")))

	user, err := db.GetUserByEmail(c.Request.Context(), h.pool, email)
	if errors.Is(err, db.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		log.Printf("get user by email: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUserByID(c.Request.Context(), h.pool, id)
	if errors.Is(err, db.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		log.Printf("get user %d: %v", id, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if req.Username != "" {
		user.Username = strings.TrimSpace(req.Username)
	}
	if req.Email != "" {
		user.Email = strings.ToLower(strings.TrimSpace(req.Email))
	}

	err = db.UpdateUser(c.Request.Context(), h.pool, user, id)
	if errors.Is(err, db.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if errors.Is(err, db.ErrAlreadyExists) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "username or email already taken"})
		return
	}
	if err != nil {
		log.Printf("update user %d: %v", id, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	err = db.DeleteUser(c.Request.Context(), h.pool, id)
	if errors.Is(err, db.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		log.Printf("delete user %d: %v", id, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Status(http.StatusNoContent)
}
