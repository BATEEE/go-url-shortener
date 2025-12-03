package handlers

import (
	"math/rand"
	"net/http"
	"simple-shortener/store"
	"time"

	"github.com/gin-gonic/gin"
)

// Dữ liệu user gửi lên
type CreateUrlRequest struct {
	OriginalUrl string `json:"original_url"`
}

// API: Tạo link rút gọn
func CreateShortLink(c *gin.Context) {
	var body CreateUrlRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// --- XỬ LÝ CONCURRENCY (TRÙNG LẶP) ---
	// Thử tạo mã tối đa 3 lần. Nếu trùng thì tạo lại cái khác.
	for i := 0; i < 3; i++ {
		shortCode := generateShortCode(6) // Tạo mã ngẫu nhiên 6 ký tự

		// Gọi Store để lưu vào DB
		err := store.SaveUrl(shortCode, body.OriginalUrl)
		
		if err == nil {
			// Thành công (không bị trùng)
			c.JSON(http.StatusOK, gin.H{
				"message":    "Rút gọn thành công",
				"short_url":  "http://localhost:8080/" + shortCode,
				"short_code": shortCode,
			})
			return
		}
		// Nếu có lỗi (err != nil) -> Vòng lặp chạy tiếp để thử mã khác
	}

	// Nếu xui quá thử 3 lần vẫn trùng
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Hệ thống đang bận, vui lòng thử lại"})
}

// API: Chuyển hướng (Redirect)
func RedirectLink(c *gin.Context) {
	code := c.Param("code") // Lấy mã trên URL

	originalUrl, exists := store.GetUrl(code)
	if exists {
		c.Redirect(http.StatusFound, originalUrl) // Chuyển hướng 302
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link không tồn tại"})
	}
}

// Hàm phụ: Sinh chuỗi ngẫu nhiên
func generateShortCode(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Khởi tạo hạt giống ngẫu nhiên
func init() {
	rand.Seed(time.Now().UnixNano())
}