package sqlcgen

import (
	"context"

	"github.com/google/uuid"
)

const createReview = `
INSERT INTO reviews (
    place_id,
    user_id,
    rating
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id, place_id, user_id, rating, created_at, updated_at
`

type CreateReviewParams struct {
	PlaceID string
	UserID  uuid.UUID
	Rating  float64
}

func (q *Queries) CreateReview(ctx context.Context, arg CreateReviewParams) (Review, error) {
	row := q.db.QueryRowContext(ctx, createReview, arg.PlaceID, arg.UserID, arg.Rating)
	var i Review
	err := row.Scan(
		&i.ID,
		&i.PlaceID,
		&i.UserID,
		&i.Rating,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
