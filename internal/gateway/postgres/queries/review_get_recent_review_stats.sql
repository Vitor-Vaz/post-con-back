-- name: GetRecentReviewStats :one
SELECT COUNT(*)::integer AS review_count, COALESCE(SUM(rating), 0)::float8 AS rating_sum
FROM (
    SELECT rating
    FROM reviews
    WHERE place_id = $1
    ORDER BY created_at DESC
    LIMIT $2
) q;
