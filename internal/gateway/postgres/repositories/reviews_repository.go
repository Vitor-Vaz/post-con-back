package repositories

import (
	"context"
	"errors"

	"github.com/lib/pq"

	"post-con-back/internal/domain"
	"post-con-back/internal/gateway/postgres/sqlcgen"
)

type ReviewsRepository struct {
	q *sqlcgen.Queries
}

func NewReviewsRepository(db sqlcgen.DBTX) *ReviewsRepository {
	return &ReviewsRepository{q: sqlcgen.New(db)}
}

func (r *ReviewsRepository) InsertReview(ctx context.Context, in domain.CreateReviewInput) (domain.Review, error) {
	row, err := r.q.CreateReview(ctx, sqlcgen.CreateReviewParams{
		PlaceID: in.PlaceID,
		UserID:  in.UserID,
		Rating:  in.Rating,
	})
	if err != nil {
		return domain.Review{}, mapInsertReviewError(err)
	}
	return domain.Review{
		ID:        row.ID,
		PlaceID:   row.PlaceID,
		UserID:    row.UserID,
		Rating:    row.Rating,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *ReviewsRepository) GetRecentReviewStats(ctx context.Context, placeID string, limit int32) (int32, float64, error) {
	row, err := r.q.GetRecentReviewStats(ctx, sqlcgen.GetRecentReviewStatsParams{
		PlaceID: placeID,
		Limit:   limit,
	})
	if err != nil {
		return 0, 0, errors.Join(domain.ErrUnexpected, err)
	}
	return row.ReviewCount, row.RatingSum, nil
}

func mapInsertReviewError(err error) error {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23514", "23502":
			return domain.ErrBadParams
		case "23505":
			return domain.ErrConflict
		}
	}
	return errors.Join(domain.ErrUnexpected, err)
}
