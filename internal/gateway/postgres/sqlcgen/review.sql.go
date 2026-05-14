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

const getRecentReviewStats = `
SELECT COUNT(*) AS review_count, COALESCE(SUM(rating), 0) AS rating_sum
FROM (
    SELECT rating
    FROM reviews
    WHERE place_id = $1
    ORDER BY created_at DESC
    LIMIT $2
) q
`

type GetRecentReviewStatsParams struct {
	PlaceID string
	Limit   int32
}

type GetRecentReviewStatsRow struct {
	ReviewCount int32   `json:"review_count"`
	RatingSum   float64 `json:"rating_sum"`
}

func (q *Queries) GetRecentReviewStats(ctx context.Context, arg GetRecentReviewStatsParams) (GetRecentReviewStatsRow, error) {
	row := q.db.QueryRowContext(ctx, getRecentReviewStats, arg.PlaceID, arg.Limit)
	var i GetRecentReviewStatsRow
	if err := row.Scan(&i.ReviewCount, &i.RatingSum); err != nil {
		return GetRecentReviewStatsRow{}, err
	}
	return i, nil
}
