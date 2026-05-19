-- name: GetStations :many
SELECT
    id,
    place_id,
    name,
    address,
    latitude,
    longitude,
    total_score,
    review_count,
    summary,
    created_at,
    updated_at
FROM station
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountStations :one
SELECT COUNT(*)::bigint AS total
FROM station;
