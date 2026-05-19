-- name: GetStationByPlaceID :one
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
WHERE place_id = $1;
