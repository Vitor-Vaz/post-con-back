package domain

import "errors"

var (
	ErrBadParams   = errors.New("bad params")
	ErrBadRequest  = errors.New("bad request")
	ErrNotFound    = errors.New("not found")
	ErrConflict    = errors.New("resource conflict")
	ErrUnexpected  = errors.New("unexpected error")
)
