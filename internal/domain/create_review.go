package domain

import "context"

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
	if err := ValidateCreateReview(in); err != nil {
		return Review{}, err
	}
	return uc.repo.InsertReview(ctx, in)
}
