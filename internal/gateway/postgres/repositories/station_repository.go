package repositories

import (
	"context"
	"errors"

	"post-con-back/internal/domain"
	"post-con-back/internal/gateway/postgres/sqlcgen"
)

type StationRepository struct {
	q *sqlcgen.Queries
}

func NewStationRepository(db sqlcgen.DBTX) *StationRepository {
	return &StationRepository{q: sqlcgen.New(db)}
}

func (r *StationRepository) UpsertStationScore(ctx context.Context, placeID string, totalScore float64, reviewCount int32) error {
	if err := r.q.UpsertStationScore(ctx, sqlcgen.UpsertStationScoreParams{
		PlaceID:     placeID,
		TotalScore:  totalScore,
		ReviewCount: reviewCount,
	}); err != nil {
		return errors.Join(domain.ErrUnexpected, err)
	}
	return nil
}
