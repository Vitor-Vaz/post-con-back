package v1_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"post-con-back/internal/domain"
)

var (
	stationsTestTime      = time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	stationsTestStationID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	stationsTestAddress   = "Rua A, 1"
	stationsTestSummary   = "Bom posto"
)

func sampleGetStationsOutput() domain.GetStationsOutput {
	return domain.GetStationsOutput{
		Stations: []domain.Station{sampleStation()},
		Page:     1,
		PageSize: domain.StationsPageSize,
		Total:    25,
	}
}

func sampleStation() domain.Station {
	return domain.Station{
		ID:          stationsTestStationID,
		PlaceID:     "ChIJx",
		Name:        "Posto Test",
		Address:     &stationsTestAddress,
		TotalScore:  4.5,
		ReviewCount: 10,
		Summary:     &stationsTestSummary,
		CreatedAt:   stationsTestTime,
		UpdatedAt:   stationsTestTime,
	}
}

func sampleStationSuccessJSON(t *testing.T) string {
	t.Helper()
	wantObj := map[string]any{
		"id":           stationsTestStationID.String(),
		"place_id":     "ChIJx",
		"name":         "Posto Test",
		"address":      stationsTestAddress,
		"total_score":  4.5,
		"review_count": float64(10),
		"summary":      stationsTestSummary,
		"created_at":   stationsTestTime.Format(time.RFC3339Nano),
		"updated_at":   stationsTestTime.Format(time.RFC3339Nano),
	}
	wantJSON, err := json.Marshal(wantObj)
	require.NoError(t, err)
	return string(wantJSON)
}

func sampleGetStationsSuccessJSON(t *testing.T) string {
	t.Helper()
	wantSuccessObj := map[string]any{
		"data": []any{
			map[string]any{
				"id":           stationsTestStationID.String(),
				"place_id":     "ChIJx",
				"name":         "Posto Test",
				"address":      stationsTestAddress,
				"total_score":  4.5,
				"review_count": float64(10),
				"summary":      stationsTestSummary,
				"created_at":   stationsTestTime.Format(time.RFC3339Nano),
				"updated_at":   stationsTestTime.Format(time.RFC3339Nano),
			},
		},
		"pagination": map[string]any{
			"page":        float64(1),
			"page_size":   float64(domain.StationsPageSize),
			"total":       float64(25),
			"total_pages": float64(3),
		},
	}
	wantSuccessJSON, err := json.Marshal(wantSuccessObj)
	require.NoError(t, err)
	return string(wantSuccessJSON)
}
