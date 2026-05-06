package repositories

import (
	"context"
	"database/sql"

	"post-con-back/internal/domain"
	"post-con-back/internal/gateway/postgres/sqlcgen"
)

type ReviewsRepository struct {
	q *sqlcgen.Queries
}

func NewReviewsRepository(db *sql.DB) *ReviewsRepository {
	return &ReviewsRepository{q: sqlcgen.New(db)}
}

func (r *ReviewsRepository) InsertReview(ctx context.Context, in domain.CreateReviewInput) (domain.Review, error) {
	row, err := r.q.CreateReview(ctx, sqlcgen.CreateReviewParams{
		PlaceID: in.PlaceID,
		UserID:  in.UserID,
		Rating:  in.Rating,
	})
	if err != nil {
		return domain.Review{}, err
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
