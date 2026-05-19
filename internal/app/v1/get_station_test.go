package v1_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "post-con-back/internal/app/v1"
	"post-con-back/internal/domain"
)

type stubStationByPlaceIDGetter struct {
	result  domain.Station
	err     error
	lastID  string
	callCnt int
}

func (s *stubStationByPlaceIDGetter) GetStationByPlaceID(ctx context.Context, placeID string) (domain.Station, error) {
	s.callCnt++
	s.lastID = placeID
	return s.result, s.err
}

func newStationRouter(listUC v1.StationsGetter, getUC v1.StationByPlaceIDGetter) *gin.Engine {
	r := gin.New()
	h := v1.NewStationsHandler(listUC, getUC)
	r.GET("/api/v1/station/:place_id", h.GetStation)
	return r
}

type getStationCase struct {
	name           string
	placeIDPath    string
	stubResult     domain.Station
	stubErr        error
	wantStatus     int
	wantErrExact   string
	wantCalls      int
	wantPlaceID    string
	wantRespJSONEq string
}

func TestGetStation(t *testing.T) {
	station := sampleStation()
	wantSuccessJSON := sampleStationSuccessJSON(t)

	tests := []getStationCase{
		{
			name:           "successful get",
			placeIDPath:    "ChIJx",
			stubResult:     station,
			wantStatus:     http.StatusOK,
			wantCalls:      1,
			wantPlaceID:    "ChIJx",
			wantRespJSONEq: wantSuccessJSON,
		},
		{
			name:           "trims place_id",
			placeIDPath:    "%20%20ChIJtrimmed%20%20",
			stubResult:     station,
			wantStatus:     http.StatusOK,
			wantCalls:      1,
			wantPlaceID:    "ChIJtrimmed",
			wantRespJSONEq: wantSuccessJSON,
		},
		{
			name:         "missing place_id",
			placeIDPath:  " ",
			wantStatus:   http.StatusBadRequest,
			wantErrExact: domain.ErrBadParams.Error(),
			wantCalls:    0,
		},
		{
			name:         "not found",
			placeIDPath:  "ChIJmissing",
			stubErr:      domain.ErrNotFound,
			wantStatus:   http.StatusNotFound,
			wantErrExact: domain.ErrNotFound.Error(),
			wantCalls:    1,
			wantPlaceID:  "ChIJmissing",
		},
		{
			name:         "unexpected error",
			placeIDPath:  "ChIJx",
			stubErr:      errors.New("db down"),
			wantStatus:   http.StatusInternalServerError,
			wantErrExact: domain.ErrUnexpected.Error(),
			wantCalls:    1,
			wantPlaceID:  "ChIJx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listStub := &stubStationsGetter{}
			getStub := &stubStationByPlaceIDGetter{result: tt.stubResult, err: tt.stubErr}
			srv := httptest.NewServer(newStationRouter(listStub, getStub))
			defer srv.Close()

			req, err := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/station/"+tt.placeIDPath, nil)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			assert.Equal(t, tt.wantCalls, getStub.callCnt)

			switch {
			case tt.wantRespJSONEq != "":
				assert.JSONEq(t, tt.wantRespJSONEq, string(bodyBytes))
			case tt.wantErrExact != "":
				var out map[string]string
				require.NoError(t, json.Unmarshal(bodyBytes, &out))
				assert.Equal(t, tt.wantErrExact, out["error"])
			}

			if tt.wantPlaceID != "" {
				assert.Equal(t, tt.wantPlaceID, getStub.lastID)
			}
		})
	}
}
