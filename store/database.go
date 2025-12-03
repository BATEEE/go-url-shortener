package store

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mô hình dữ liệu (Database Schema)
type ShortLink struct {
	ID          string `gorm:"primaryKey"` // Mã rút gọn (VD: abc) - Là khóa chính
	OriginalUrl string // Link gốc
	Clicks      int    // Số lượt click
}

var db *gorm.DB

// 1. Khởi động Database (Tự tạo file urls.db)
func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Lỗi kết nối database: ", err)
	}

	// Tự động tạo bảng nếu chưa có
	db.AutoMigrate(&ShortLink{})
}

// 2. Hàm Lưu Link (Trả về error để biết nếu bị trùng ID)
func SaveUrl(shortCode string, originalUrl string) error {
	link := ShortLink{
		ID:          shortCode,
		OriginalUrl: originalUrl,
		Clicks:      0,
	}
	// Nếu trùng ID, hàm Create sẽ trả về error
	result := db.Create(&link)
	return result.Error
}

// 3. Hàm Lấy Link
func GetUrl(shortCode string) (string, bool) {
	var link ShortLink
	// Tìm link theo ID
	result := db.First(&link, "id = ?", shortCode)

	if result.Error != nil {
		return "", false // Không tìm thấy
	}

	// Tăng lượt click (Bonus feature)
	link.Clicks++
	db.Save(&link)

	return link.OriginalUrl, true
}