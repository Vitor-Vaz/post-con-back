package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID
	PlaceID   string
	UserID    uuid.UUID
	Rating    float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateReviewInput struct {
	PlaceID string
	UserID  uuid.UUID
	Rating  float64
}

func ValidateCreateReview(in CreateReviewInput) error {
	if strings.TrimSpace(in.PlaceID) == "" {
		return ErrEmptyPlaceID
	}
	if in.UserID == uuid.Nil {
		return ErrInvalidUserID
	}
	if in.Rating < 1 || in.Rating > 5 {
		return ErrInvalidReviewRating
	}
	return nil
}
