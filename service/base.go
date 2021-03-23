package service

import (
	"fmt"
	"goblog/repository"
)

type Service struct {
	repository *repository.Repository
}

func NewService() *Service {
	return &Service{
		repository: repository.NewRepository(),
	}
}

var (
	// ErrNoRow is usually happen in get list.
	ErrNoRow = fmt.Errorf("no row has been selected")
)

// define Error used in display system.

type DispError struct {
	Type DispErrorType
	Err  error
}

type DispErrorType uint8

const (
	ErrorRequest DispErrorType = iota
	ErrorServer
)

func NewDispError(t DispErrorType, err error) *DispError {
	return &DispError{Type: t, Err: err}
}

func (e *DispError) GetErrorType() DispErrorType {
	return e.Type
}

func (e *DispError) GetError() error {
	return e.Err
}
