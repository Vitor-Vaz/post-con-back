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

type GetStationsOutput struct {
	Stations []Station
	Page     int
	PageSize int
	Total    int64
}

type StationsGetterRepository interface {
	GetStations(ctx context.Context, limit int32, offset int32) ([]Station, error)
	CountStations(ctx context.Context) (int64, error)
}

type GetStationsUseCase struct {
	repo StationsGetterRepository
}

func NewGetStationsUseCase(repo StationsGetterRepository) *GetStationsUseCase {
	return &GetStationsUseCase{repo: repo}
}

func (uc *GetStationsUseCase) GetStations(ctx context.Context, page int) (GetStationsOutput, error) {
	offset := int32((page - 1) * StationsPageSize)
	stations, err := uc.repo.GetStations(ctx, StationsPageSize, offset)
	if err != nil {
		return GetStationsOutput{}, err
	}
	total, err := uc.repo.CountStations(ctx)
	if err != nil {
		return GetStationsOutput{}, err
	}
	return GetStationsOutput{
		Stations: stations,
		Page:     page,
		PageSize: StationsPageSize,
		Total:    total,
	}, nil
}
