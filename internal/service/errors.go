package service

import "errors"

var (
	ErrHumanNotFound = errors.New("human not found")
	ErrInvalidAge = errors.New("invalid age; empty")
	ErrInvalidGender = errors.New("invalid gender; empty")
	ErrInvalidNationality = errors.New("invalid nationality; empty")
)