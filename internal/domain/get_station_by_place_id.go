package domain

import "context"

type StationByPlaceIDGetterRepository interface {
	GetStationByPlaceID(ctx context.Context, placeID string) (Station, error)
}

type GetStationByPlaceIDUseCase struct {
	repo StationByPlaceIDGetterRepository
}

func NewGetStationByPlaceIDUseCase(repo StationByPlaceIDGetterRepository) *GetStationByPlaceIDUseCase {
	return &GetStationByPlaceIDUseCase{repo: repo}
}

func (uc *GetStationByPlaceIDUseCase) GetStationByPlaceID(ctx context.Context, placeID string) (Station, error) {
	return uc.repo.GetStationByPlaceID(ctx, placeID)
}
