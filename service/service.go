package service

import (
	"errors"
	"log"
	"net/url"
	"simple-shortener/domain"
	"simple-shortener/utils"
)

var (
	ErrInvalidURL  = errors.New("URL is invalid")
	ErrInvalidCode = errors.New("Shortener code is invalid")
	ErrCodeExists  = errors.New("Exists shortener code")
	ErrGenFailed   = errors.New("System is busy, try again later")
	ErrNotFound    = errors.New("URL not found")
	ErrEmailExists = errors.New("Email already exists")
)

type Store interface {
	CreateShortLink(link *domain.Link) error
	GetByShortCode(code string) (*domain.Link, error)
	CreateUser(u *domain.User) error
	IncrementClicks(code string) error
}

type Service struct {
	store Store
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) CreateShort(userID uint64, shortCode string, originalURL string) (string, error) {
	// 1. Validate URL
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return "", ErrInvalidURL
	}

	// CASE 1: User Input
	if shortCode != "" {
		if !isValidShortCode(shortCode) {
			return "", ErrInvalidCode
		}
		link := &domain.Link{
			UserID:      userID,
			ShortCode:   shortCode,
			OriginalUrl: originalURL,
			Clicks:      0,
		}
		if err := s.store.CreateShortLink(link); err != nil {
			return "", err
		}
		return shortCode, nil
	}

	// CASE 2: Random code
	for i := 0; i < 5; i++ {
		newCode := utils.GenerateRandomString(6)

		link := &domain.Link{
			UserID:      userID,
			ShortCode:   newCode,
			OriginalUrl: originalURL,
			Clicks:      0,
		}

		err := s.store.CreateShortLink(link)

		if err == nil {
			return newCode, nil
		}

		if err == ErrCodeExists {
			continue
		}

		return "", err
	}

	return "", ErrGenFailed
}

func (s *Service) GetOriginalAndIncrement(code string) (string, error) {
	link, err := s.store.GetByShortCode(code)
	if err != nil {
		return "", err
	}

	go func() {
		if err := s.store.IncrementClicks(code); err != nil {
			log.Printf("Failed to increment clicks for %s: %v", code, err)
		}
	}()

	return link.OriginalUrl, nil
}

func isValidShortCode(code string) bool {
	if len(code) < 3 || len(code) > 32 {
		return false
	}
	for i := 0; i < len(code); i++ {
		c := code[i]
		if !(c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c == '-' || c == '_') {
			return false
		}
	}
	return true
}

// ---------------------------USER--------------------------
func (s *Service) CreateUser(email string) (*domain.User, error) {
	user := &domain.User{
		Email: email,
	}

	if err := s.store.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
