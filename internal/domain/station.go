package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const StationsPageSize = 10

type Station struct {
	ID          uuid.UUID
	PlaceID     string
	Name        string
	Address     *string
	Latitude    *float64
	Longitude   *float64
	TotalScore  float64
	ReviewCount int32
	Summary     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ListStationsOutput struct {
	Stations []Station
	Page     int
	PageSize int
	Total    int64
}

type StationsListerRepository interface {
	ListStations(ctx context.Context, limit int32, offset int32) ([]Station, error)
	CountStations(ctx context.Context) (int64, error)
}

type ListStationsUseCase struct {
	repo StationsListerRepository
}

func NewListStationsUseCase(repo StationsListerRepository) *ListStationsUseCase {
	return &ListStationsUseCase{repo: repo}
}

func (uc *ListStationsUseCase) ListStations(ctx context.Context, page int) (ListStationsOutput, error) {
	offset := int32((page - 1) * StationsPageSize)
	stations, err := uc.repo.ListStations(ctx, StationsPageSize, offset)
	if err != nil {
		return ListStationsOutput{}, err
	}
	total, err := uc.repo.CountStations(ctx)
	if err != nil {
		return ListStationsOutput{}, err
	}
	return ListStationsOutput{
		Stations: stations,
		Page:     page,
		PageSize: StationsPageSize,
		Total:    total,
	}, nil
}
