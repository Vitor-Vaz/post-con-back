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
}

type StationRepository interface {
	UpsertStationScore(ctx context.Context, placeID string, totalScore float64, reviewCount int32) error
}

type ReviewCreatorUseCase struct {
	reviewRepo  ReviewCreatorRepository
	stationRepo StationRepository
}

func NewReviewCreatorUseCase(reviewRepo ReviewCreatorRepository, stationRepo StationRepository) *ReviewCreatorUseCase {
	return &ReviewCreatorUseCase{reviewRepo: reviewRepo, stationRepo: stationRepo}
}

func (uc *ReviewCreatorUseCase) CreateReview(ctx context.Context, in CreateReviewInput) (Review, error) {
	review, err := uc.reviewRepo.InsertReview(ctx, in)
	if err != nil {
		return Review{}, err
	}

	reviewCount, ratingSum, err := uc.reviewRepo.GetRecentReviewStats(ctx, in.PlaceID, reviewAverageWindow)
	if err != nil {
		return Review{}, err
	}
	if reviewCount == 0 {
		reviewCount = 1
		ratingSum = review.Rating
	}

	average := ratingSum / float64(reviewCount)
	if err := uc.stationRepo.UpsertStationScore(ctx, in.PlaceID, average, reviewCount); err != nil {
		return Review{}, err
	}

	return review, nil
}
