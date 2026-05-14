package sqlcgen

import (
	"context"
)

const upsertStationScore = `
INSERT INTO station (
    place_id,
    name,
    total_score,
    review_count,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    now()
)
ON CONFLICT (place_id) DO UPDATE SET
    total_score = EXCLUDED.total_score,
    review_count = EXCLUDED.review_count,
    updated_at = now();
`

type UpsertStationScoreParams struct {
	PlaceID     string
	TotalScore  float64
	ReviewCount int32
}

func (q *Queries) UpsertStationScore(ctx context.Context, arg UpsertStationScoreParams) error {
	if _, err := q.db.ExecContext(ctx, upsertStationScore, arg.PlaceID, arg.PlaceID, arg.TotalScore, arg.ReviewCount); err != nil {
		return err
	}
	return nil
}
