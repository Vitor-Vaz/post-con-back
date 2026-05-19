package v1_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "post-con-back/internal/app/v1"
	"post-con-back/internal/domain"
)

const stationsAPIPath = "/api/v1/stations"

type stubStationsLister struct {
	result  domain.ListStationsOutput
	err     error
	lastPage int
	callCnt int
}

func (s *stubStationsLister) ListStations(ctx context.Context, page int) (domain.ListStationsOutput, error) {
	s.callCnt++
	s.lastPage = page
	return s.result, s.err
}

func newStationsRouter(uc v1.StationsLister) *gin.Engine {
	r := gin.New()
	h := v1.NewStationsHandler(uc)
	r.GET(stationsAPIPath, h.GetStations)
	return r
}

type getStationsCase struct {
	name           string
	query          string
	stubResult     domain.ListStationsOutput
	stubErr        error
	wantStatus     int
	wantErrExact   string
	wantCalls      int
	wantPage       int
	wantRespJSONEq string
}

func TestGetStations(t *testing.T) {
	ts := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	stationID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	address := "Rua A, 1"
	summary := "Bom posto"

	successOut := domain.ListStationsOutput{
		Stations: []domain.Station{
			{
				ID:          stationID,
				PlaceID:     "ChIJx",
				Name:        "Posto Test",
				Address:     &address,
				TotalScore:  4.5,
				ReviewCount: 10,
				Summary:     &summary,
				CreatedAt:   ts,
				UpdatedAt:   ts,
			},
		},
		Page:     1,
		PageSize: domain.StationsPageSize,
		Total:    25,
	}
	wantSuccessObj := map[string]any{
		"data": []any{
			map[string]any{
				"id":           stationID.String(),
				"place_id":     "ChIJx",
				"name":         "Posto Test",
				"address":      address,
				"total_score":  4.5,
				"review_count": float64(10),
				"summary":      summary,
				"created_at":   ts.Format(time.RFC3339Nano),
				"updated_at":   ts.Format(time.RFC3339Nano),
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

	tests := []getStationsCase{
		{
			name:           "successful list default page",
			stubResult:     successOut,
			wantStatus:     http.StatusOK,
			wantCalls:      1,
			wantPage:       1,
			wantRespJSONEq: string(wantSuccessJSON),
		},
		{
			name:           "successful list explicit page",
			query:          "?page=2",
			stubResult:     successOut,
			wantStatus:     http.StatusOK,
			wantCalls:      1,
			wantPage:       2,
			wantRespJSONEq: string(wantSuccessJSON),
		},
		{
			name:         "invalid page zero",
			query:        "?page=0",
			wantStatus:   http.StatusBadRequest,
			wantErrExact: domain.ErrBadParams.Error(),
			wantCalls:    0,
		},
		{
			name:         "invalid page negative",
			query:        "?page=-1",
			wantStatus:   http.StatusBadRequest,
			wantErrExact: domain.ErrBadParams.Error(),
			wantCalls:    0,
		},
		{
			name:         "invalid page not a number",
			query:        "?page=abc",
			wantStatus:   http.StatusBadRequest,
			wantErrExact: domain.ErrBadParams.Error(),
			wantCalls:    0,
		},
		{
			name:         "use case unexpected",
			stubErr:      errors.New("db unavailable"),
			wantStatus:   http.StatusInternalServerError,
			wantErrExact: domain.ErrUnexpected.Error(),
			wantCalls:    1,
			wantPage:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &stubStationsLister{result: tt.stubResult, err: tt.stubErr}
			srv := httptest.NewServer(newStationsRouter(stub))
			defer srv.Close()

			req, err := http.NewRequest(http.MethodGet, srv.URL+stationsAPIPath+tt.query, nil)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			assert.Equal(t, tt.wantCalls, stub.callCnt)

			switch {
			case tt.wantRespJSONEq != "":
				assert.JSONEq(t, tt.wantRespJSONEq, string(bodyBytes))
			case tt.wantErrExact != "":
				var out map[string]string
				require.NoError(t, json.Unmarshal(bodyBytes, &out))
				assert.Equal(t, tt.wantErrExact, out["error"])
			}

			if tt.wantPage != 0 {
				assert.Equal(t, tt.wantPage, stub.lastPage)
			}
		})
	}
}
