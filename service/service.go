package service

import (
	"errors"
	"log"
	"net/url"
	"simple-shortener/domain"
	"simple-shortener/utils"
)

var (
	ErrInvalidURL  = domain.ErrInvalidURL
	ErrInvalidCode = domain.ErrInvalidCode
	ErrCodeExists  = domain.ErrCodeExists
	ErrGenFailed   = domain.ErrGenFailed
	ErrNotFound    = domain.ErrNotFound
	ErrEmailExists = domain.ErrEmailExists
)

type Store interface {
	CreateShortLink(link *domain.Link) error
	GetByShortCode(code string) (*domain.Link, error)
	GetByOriginalURL(userID uint64, originalURL string) (*domain.Link, error)
	IncrementClicks(code string) error
	CreateUser(u *domain.User) error
	GetUserByID(userID uint64) (*domain.User, error)
	GetLinksByUserID(userID uint64) ([]*domain.Link, error)
}

type Service struct {
	store Store
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) CreateShort(userID uint64, shortCode string, originalURL string) (string, error) {
	// Check if user exists
	_, err := s.store.GetUserByID(userID)
	if err != nil {
		if err == ErrNotFound {
			return "", errors.New("user not found")
		}
		return "", err
	}

	// Validate URL
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return "", ErrInvalidURL
	}

	// Check short_code provided
	if shortCode != "" {
		if !isValidShortCode(shortCode) {
			return "", ErrInvalidCode
		}

		// Check if this user already created a link for this URL
		existingLink, err := s.store.GetByOriginalURL(userID, originalURL)
		if err == nil && existingLink != nil {
			return "", errors.New("URL already shortened: " + existingLink.ShortCode)
		}

		existingByCode, err := s.store.GetByShortCode(shortCode)
		if err == nil && existingByCode != nil {
			if existingByCode.UserID == userID && existingByCode.OriginalUrl == originalURL {
				return existingByCode.ShortCode, nil
			}
			return "", ErrCodeExists
		}
		if err != nil && err != ErrNotFound {
			return "", err
		}
	}

	// Check if this user already created a link for this URL
	existingLink, err := s.store.GetByOriginalURL(userID, originalURL)
	if err == nil && existingLink != nil {
		return "", errors.New("URL already shortened: " + existingLink.ShortCode)
	}

	if err != nil && err != ErrNotFound {
		return "", err
	}

	if shortCode != "" {
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

	// Random code generation
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

func (s *Service) GetLinkInfo(code string) (*domain.Link, error) {
	return s.store.GetByShortCode(code)
}

func (s *Service) GetUserLinks(userID uint64) ([]*domain.Link, error) {
	// Check if user exists
	_, err := s.store.GetUserByID(userID)
	if err != nil {
		if err == ErrNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.store.GetLinksByUserID(userID)
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
