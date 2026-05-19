package repositories

import (
	"context"
	"database/sql"
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
		Name:        placeID,
		TotalScore:  totalScore,
		ReviewCount: reviewCount,
	}); err != nil {
		return errors.Join(domain.ErrUnexpected, err)
	}
	return nil
}

func (r *StationRepository) ListStations(ctx context.Context, limit int32, offset int32) ([]domain.Station, error) {
	rows, err := r.q.ListStations(ctx, sqlcgen.ListStationsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, errors.Join(domain.ErrUnexpected, err)
	}
	stations := make([]domain.Station, 0, len(rows))
	for _, row := range rows {
		stations = append(stations, mapStationRow(row))
	}
	return stations, nil
}

func (r *StationRepository) CountStations(ctx context.Context) (int64, error) {
	total, err := r.q.CountStations(ctx)
	if err != nil {
		return 0, errors.Join(domain.ErrUnexpected, err)
	}
	return total, nil
}

func mapStationRow(row sqlcgen.Station) domain.Station {
	return domain.Station{
		ID:          row.ID,
		PlaceID:     row.PlaceID,
		Name:        row.Name,
		Address:     nullStringToPtr(row.Address),
		Latitude:    nullFloat64ToPtr(row.Latitude),
		Longitude:   nullFloat64ToPtr(row.Longitude),
		TotalScore:  row.TotalScore,
		ReviewCount: row.ReviewCount,
		Summary:     nullStringToPtr(row.Summary),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func nullStringToPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func nullFloat64ToPtr(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	f := v.Float64
	return &f
}
