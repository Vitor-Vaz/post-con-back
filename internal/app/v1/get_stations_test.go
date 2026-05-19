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

type stubStationsGetter struct {
	result   domain.GetStationsOutput
	err      error
	lastPage int
	callCnt  int
}

func (s *stubStationsGetter) GetStations(ctx context.Context, page int) (domain.GetStationsOutput, error) {
	s.callCnt++
	s.lastPage = page
	return s.result, s.err
}

func newStationsRouter(uc v1.StationsGetter) *gin.Engine {
	r := gin.New()
	h := v1.NewStationsHandler(uc)
	r.GET("/api/v1/stations", h.GetStations)
	return r
}

type getStationsCase struct {
	name           string
	query          string
	stubResult     domain.GetStationsOutput
	stubErr        error
	wantStatus     int
	wantErrExact   string
	wantCalls      int
	wantPage       int
	wantRespJSONEq string
}

func TestGetStations(t *testing.T) {
	successOut := sampleGetStationsOutput()
	wantSuccessJSON := sampleGetStationsSuccessJSON(t)

	tests := []getStationsCase{
		{
			name:           "successful list default page",
			stubResult:     successOut,
			wantStatus:     http.StatusOK,
			wantCalls:      1,
			wantPage:       1,
			wantRespJSONEq: wantSuccessJSON,
		},
		{
			name:           "successful list explicit page",
			query:          "?page=2",
			stubResult:     successOut,
			wantStatus:     http.StatusOK,
			wantCalls:      1,
			wantPage:       2,
			wantRespJSONEq: wantSuccessJSON,
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
			stub := &stubStationsGetter{result: tt.stubResult, err: tt.stubErr}
			srv := httptest.NewServer(newStationsRouter(stub))
			defer srv.Close()

			req, err := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/stations"+tt.query, nil)
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
