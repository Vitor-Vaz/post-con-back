package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateCreateReview(t *testing.T) {
	uid := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	t.Run("accepts valid input", func(t *testing.T) {
		err := ValidateCreateReview(CreateReviewInput{
			PlaceID: "ChIJN1t_tDeuEmsRUsoyG83frY4",
			UserID:  uid,
			Rating:  4.5,
		})
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})
	t.Run("rejects empty place id", func(t *testing.T) {
		err := ValidateCreateReview(CreateReviewInput{PlaceID: "  ", UserID: uid, Rating: 3})
		if err != ErrEmptyPlaceID {
			t.Fatalf("expected ErrEmptyPlaceID, got %v", err)
		}
	})
	t.Run("rejects nil user id", func(t *testing.T) {
		err := ValidateCreateReview(CreateReviewInput{PlaceID: "x", UserID: uuid.Nil, Rating: 3})
		if err != ErrInvalidUserID {
			t.Fatalf("expected ErrInvalidUserID, got %v", err)
		}
	})
	t.Run("rejects rating below range", func(t *testing.T) {
		err := ValidateCreateReview(CreateReviewInput{PlaceID: "x", UserID: uid, Rating: 0.9})
		if err != ErrInvalidReviewRating {
			t.Fatalf("expected ErrInvalidReviewRating, got %v", err)
		}
	})
	t.Run("rejects rating above range", func(t *testing.T) {
		err := ValidateCreateReview(CreateReviewInput{PlaceID: "x", UserID: uid, Rating: 5.1})
		if err != ErrInvalidReviewRating {
			t.Fatalf("expected ErrInvalidReviewRating, got %v", err)
		}
	})
}
