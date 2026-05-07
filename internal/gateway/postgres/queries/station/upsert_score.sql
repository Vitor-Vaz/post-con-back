-- name: UpsertStationScore :exec
INSERT INTO station (
    place_id,
    total_score,
    review_count,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    now()
)
ON CONFLICT (place_id) DO UPDATE SET
    total_score = EXCLUDED.total_score,
    review_count = EXCLUDED.review_count,
    updated_at = now();
