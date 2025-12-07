package repository

import (
	"errors"
	"fmt"
	"simple-shortener/domain"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(path string) (*GormStore, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Link{}); err != nil {
		return nil, err
	}

	return &GormStore{db: db}, nil
}

func (s *GormStore) CreateUser(u *domain.User) error {
	if u == nil {
		return fmt.Errorf("nil user")
	}
	if err := s.db.Create(u).Error; err != nil {
		if isUniqueConstraintError(err) {
			return domain.ErrEmailExists
		}
		return err
	}
	return nil
}

func (s *GormStore) GetUserByID(userID uint64) (*domain.User, error) {
	var user domain.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *GormStore) CreateShortLink(link *domain.Link) error {
	if err := s.db.Create(link).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "constraint") {
			return domain.ErrCodeExists
		}
		return err
	}
	return nil
}

func (s *GormStore) GetByShortCode(code string) (*domain.Link, error) {
	var l domain.Link
	if err := s.db.First(&l, "short_code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (s *GormStore) IncrementClicks(code string) error {
	res := s.db.Model(&domain.Link{}).
		Where("short_code = ?", code).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()
	return msg != "" && (strings.Contains(msg, "UNIQUE constraint failed") || strings.Contains(msg, "UNIQUE constraint"))
}

func (s *GormStore) GetLinksByUserID(userID uint64) ([]*domain.Link, error) {
	var links []*domain.Link
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (s *GormStore) GetByOriginalURL(userID uint64, originalURL string) (*domain.Link, error) {
	var l domain.Link
	if err := s.db.First(&l, "user_id = ? AND original_url = ?", userID, originalURL).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}
