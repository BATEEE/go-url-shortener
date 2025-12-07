package domain

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidURL  = errors.New("URL is invalid")
	ErrInvalidCode = errors.New("shortener code is invalid")
	ErrCodeExists  = errors.New("shortener code exists")
	ErrGenFailed   = errors.New("system is busy, try again later")
	ErrEmailExists = errors.New("email already exists")
)
