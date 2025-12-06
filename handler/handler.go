package handler

import (
	"errors"
	"net/http"
	"simple-shortener/domain"
	"simple-shortener/domain/dto"
	"simple-shortener/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateShort(userID uint64, shortCode string, originalURL string) (string, error)
	GetOriginalAndIncrement(code string) (string, error)
	CreateUser(email string) (*domain.User, error)
}

type HTTPHandler struct {
	svc Service
	r   *gin.Engine
}

func NewHandler(s Service) *HTTPHandler {
	r := gin.Default()
	h := &HTTPHandler{svc: s, r: r}

	r.POST("/users", h.createUser)

	r.POST("/users/:user_id/shorten", h.createShort)
	r.GET("/:code", h.getAndRedirect)

	return h
}

func (h *HTTPHandler) Run(addr string) error {
	return h.r.Run(addr)
}

func (h *HTTPHandler) createUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	user, err := h.svc.CreateUser(req.Email)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists or system error"})
		return
	}

	resp := dto.CreateUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *HTTPHandler) createShort(c *gin.Context) {
	var req dto.CreateShortLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	uidStr := c.Param("user_id")
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	code, err := h.svc.CreateShort(uid, req.ShortCode, req.URL)
	if err != nil {
		if errors.Is(err, service.ErrCodeExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Short code already exists"})
			return
		}
		if errors.Is(err, service.ErrInvalidURL) || errors.Is(err, service.ErrInvalidCode) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	resp := dto.CreateShortLinkResponse{
		ShortCode:   code,
		OriginalUrl: req.URL,
		UserID:      uid,
		Clicks:      0,
		CreatedAt:   time.Now(),
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *HTTPHandler) getAndRedirect(c *gin.Context) {
	code := c.Param("code")

	orig, err := h.svc.GetOriginalAndIncrement(code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Redirect thật (302 Found)
	c.Redirect(http.StatusFound, orig)
}
