// Package repository contains all application repositories.
package repository

import "errors"

var (
	// ErrNotFound represents a not found error.
	ErrNotFound = errors.New("not found")
	// ErrInvalidID represents an invalid id error.
	ErrInvalidID = errors.New("id provided was not a valid UUID")
)
