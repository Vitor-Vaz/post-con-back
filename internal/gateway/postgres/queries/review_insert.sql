-- name: CreateReview :one
INSERT INTO reviews (
    place_id,
    user_id,
    rating
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id, place_id, user_id, rating, created_at, updated_at;
