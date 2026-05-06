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

type ReviewCreatorRepository interface {
	InsertReview(ctx context.Context, in CreateReviewInput) (Review, error)
}

type ReviewCreatorUseCase struct {
	repo ReviewCreatorRepository
}

func NewReviewCreatorUseCase(repo ReviewCreatorRepository) *ReviewCreatorUseCase {
	return &ReviewCreatorUseCase{repo: repo}
}

func (uc *ReviewCreatorUseCase) CreateReview(ctx context.Context, in CreateReviewInput) (Review, error) {
	return uc.repo.InsertReview(ctx, in)
}
