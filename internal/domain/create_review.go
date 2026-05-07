package domain

import (
	"context"
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

const reviewAverageWindow = 100

type ReviewCreatorRepository interface {
	InsertReview(ctx context.Context, in CreateReviewInput) (Review, error)
	GetRecentReviewStats(ctx context.Context, placeID string, limit int32) (int32, float64, error)
	UpsertStationScore(ctx context.Context, placeID string, totalScore float64, reviewCount int32) error
}

type ReviewCreatorUseCase struct {
	repo ReviewCreatorRepository
}

func NewReviewCreatorUseCase(repo ReviewCreatorRepository) *ReviewCreatorUseCase {
	return &ReviewCreatorUseCase{repo: repo}
}

func (uc *ReviewCreatorUseCase) CreateReview(ctx context.Context, in CreateReviewInput) (Review, error) {
	review, err := uc.repo.InsertReview(ctx, in)
	if err != nil {
		return Review{}, err
	}

	reviewCount, ratingSum, err := uc.repo.GetRecentReviewStats(ctx, in.PlaceID, reviewAverageWindow)
	if err != nil {
		return Review{}, err
	}
	if reviewCount == 0 {
		return review, nil
	}

	average := ratingSum / float64(reviewCount)
	if err := uc.repo.UpsertStationScore(ctx, in.PlaceID, average, reviewCount); err != nil {
		return Review{}, err
	}

	return review, nil
}
