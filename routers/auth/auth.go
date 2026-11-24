package auth

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"taskforge/db"
	"taskforge/models"
)

const tokenTTL = time.Hour

type Handler struct {
	pool   *pgxpool.Pool
	jwtKey []byte
}

func NewHandler(pool *pgxpool.Pool, jwtKey []byte) *Handler {
	return &Handler{pool: pool, jwtKey: jwtKey}
}

func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("hash password: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	user := models.User{
		Username:     strings.TrimSpace(req.Username),
		Email:        normalizeEmail(req.Email),
		Role:         "user",
		PasswordHash: string(hash),
	}

	err = db.CreateUser(c.Request.Context(), h.pool, &user)
	if errors.Is(err, db.ErrAlreadyExists) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "username or email already taken"})
		return
	}
	if err != nil {
		log.Printf("create user: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUserByEmail(c.Request.Context(), h.pool, normalizeEmail(req.Email))
	if errors.Is(err, db.ErrNotFound) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err != nil {
		log.Printf("get user by email: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := h.newToken(user)
	if err != nil {
		log.Printf("sign token: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, models.Token{TokenType: "Bearer", AccessToken: token})
}

func (h *Handler) newToken(user models.User) (string, error) {
	now := time.Now()
	claims := models.Claims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(h.jwtKey)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
