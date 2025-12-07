package handler

import (
	"errors"
	"net/http"
	"simple-shortener/domain"
	"simple-shortener/domain/dto"
	"simple-shortener/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateShort(userID uint64, shortCode string, originalURL string) (string, error)
	GetOriginalAndIncrement(code string) (string, error)
	CreateUser(email string) (*domain.User, error)
	GetLinkInfo(code string) (*domain.Link, error)
	GetUserLinks(userID uint64) ([]*domain.Link, error)
}

type HTTPHandler struct {
	svc     Service
	r       *gin.Engine
	baseURL string
}

func NewHandler(s Service, baseURL string) *HTTPHandler {
	r := gin.Default()
	h := &HTTPHandler{svc: s, r: r, baseURL: baseURL}

	r.POST("/users", h.createUser)
	r.POST("/users/:user_id/shorten", h.createShort)

	r.GET("/links/:code/info", h.getLinkInfo)
	r.GET("/users/:user_id/links", h.getUserLinks)

	r.GET("/:code", h.getAndRedirect)

	return h
}

func (h *HTTPHandler) Run(addr string) error {
	return h.r.Run(addr)
}

func (h *HTTPHandler) getLinkInfo(c *gin.Context) {
	code := c.Param("code")
	link, err := h.svc.GetLinkInfo(code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, link)
}

func (h *HTTPHandler) getUserLinks(c *gin.Context) {
	uidStr := c.Param("user_id")
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	links, err := h.svc.GetUserLinks(uid)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, links)
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
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if errors.Is(err, service.ErrCodeExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Short code already exists"})
			return
		}
		// Check for duplicate URL error
		if strings.Contains(err.Error(), "URL already shortened") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
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
		ShortURL:    h.baseURL + "/" + code, // Computed from baseURL
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
