package domain

import "errors"

var (
	ErrEmptyPlaceID        = errors.New("empty place id")
	ErrInvalidReviewRating = errors.New("invalid review rating")
	ErrInvalidUserID       = errors.New("invalid user id")
)
