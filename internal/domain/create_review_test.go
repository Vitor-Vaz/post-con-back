package domain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubReviewRepository struct {
	insertReviewResult Review
	insertReviewErr    error
	statsCount         int32
	statsSum           float64
	statsErr           error
	upsertErr          error
	upsertScore        float64
	upsertCount        int32
}

func (s *stubReviewRepository) InsertReview(ctx context.Context, in CreateReviewInput) (Review, error) {
	return s.insertReviewResult, s.insertReviewErr
}

func (s *stubReviewRepository) GetRecentReviewStats(ctx context.Context, placeID string, limit int32) (int32, float64, error) {
	return s.statsCount, s.statsSum, s.statsErr
}

func (s *stubReviewRepository) UpsertStationScore(ctx context.Context, placeID string, totalScore float64, reviewCount int32) error {
	s.upsertScore = totalScore
	s.upsertCount = reviewCount
	return s.upsertErr
}

func TestReviewCreatorUseCase_CreateReview_UpdatesStationAverage(t *testing.T) {
	uid := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	repo := &stubReviewRepository{
		insertReviewResult: Review{ID: id, PlaceID: "ChIJstation", UserID: uid, Rating: 3.0, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		statsCount:         3,
		statsSum:           12.0,
	}

	uc := NewReviewCreatorUseCase(repo)

	review, err := uc.CreateReview(context.Background(), CreateReviewInput{PlaceID: "ChIJstation", UserID: uid, Rating: 3.0})
	require.NoError(t, err)
	assert.Equal(t, id, review.ID)
	assert.Equal(t, 4.0, repo.upsertScore)
	assert.Equal(t, int32(3), repo.upsertCount)
}

func TestReviewCreatorUseCase_CreateReview_BadInsertReturnsError(t *testing.T) {
	repo := &stubReviewRepository{insertReviewErr: errors.New("insert failed")}
	uc := NewReviewCreatorUseCase(repo)

	_, err := uc.CreateReview(context.Background(), CreateReviewInput{PlaceID: "ChIJstation", UserID: uuid.New(), Rating: 4.0})
	require.Error(t, err)
}
