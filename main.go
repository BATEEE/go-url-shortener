package main

import (
	"log"
	"simple-shortener/handlers"
	"simple-shortener/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Kết nối Database trước khi chạy server
	store.InitDB()

	// 2. Tạo router
	r := gin.Default()

	// 3. Đăng ký các đường dẫn (API)
	r.POST("/shorten", handlers.CreateShortLink) // Tạo link
	r.GET("/:code", handlers.RedirectLink)       // Chuyển hướng

	// 4. Chạy server tại cổng 8080
	r.Run(":8080")

	log.Println("------Server is starting!!!")
}
