package v1_test

import (
	"bytes"
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

const reviewAPIPath = "/api/v1/review"

func init() {
	gin.SetMode(gin.TestMode)
}

type stubReviewCreator struct {
	result  domain.Review
	err     error
	lastIn  domain.CreateReviewInput
	callCnt int
}

func (s *stubReviewCreator) CreateReview(ctx context.Context, in domain.CreateReviewInput) (domain.Review, error) {
	s.callCnt++
	s.lastIn = in
	return s.result, s.err
}

func newReviewRouter(uc v1.ReviewCreator) *gin.Engine {
	r := gin.New()
	h := v1.NewReviewHandler(uc)
	r.POST(reviewAPIPath, h.CreateReview)
	return r
}

type createReviewCase struct {
	name               string
	body               string
	stubResult         domain.Review
	stubErr            error
	wantStatus         int
	wantErrExact       string
	wantCalls          int
	wantRespJSONEq     string
	wantTrimmedPlaceID string
}

func TestCreateReview(t *testing.T) {
	uid := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	ts := time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC)
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	successReview := domain.Review{
		ID:        id,
		PlaceID:   "ChIJtrimmed",
		UserID:    uid,
		Rating:    4.5,
		CreatedAt: ts,
		UpdatedAt: ts,
	}
	wantSuccessObj := map[string]any{
		"id":         id.String(),
		"place_id":   "ChIJtrimmed",
		"user_id":    uid.String(),
		"rating":     4.5,
		"created_at": ts.Format(time.RFC3339Nano),
		"updated_at": ts.Format(time.RFC3339Nano),
	}
	wantSuccessJSON, err := json.Marshal(wantSuccessObj)
	require.NoError(t, err)

	validBody := `{"place_id":"ChIJx","user_id":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","rating":3}`

	tests := []createReviewCase{
		{
			name:               "successful creation",
			body:               `{"place_id":"  ChIJtrimmed  ","user_id":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","rating":4.5}`,
			stubResult:         successReview,
			wantStatus:         http.StatusCreated,
			wantCalls:          1,
			wantRespJSONEq:     string(wantSuccessJSON),
			wantTrimmedPlaceID: "ChIJtrimmed",
		},
		{
			name:         "invalid json",
			body:         `{`,
			wantStatus:   http.StatusBadRequest,
			wantErrExact: domain.ErrBadRequest.Error(),
			wantCalls:    0,
		},
		{
			name:       "missing place_id",
			body:       `{"user_id":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","rating":3}`,
			wantStatus: http.StatusBadRequest,
			wantCalls:  0,
		},
		{
			name:         "place_id whitespace only",
			body:         `{"place_id":"   ","user_id":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","rating":3}`,
			wantStatus:   http.StatusBadRequest,
			wantErrExact: domain.ErrBadParams.Error(),
			wantCalls:    0,
		},
		{
			name:       "invalid user_id format",
			body:       `{"place_id":"ChIJx","user_id":"not-a-uuid","rating":3}`,
			wantStatus: http.StatusBadRequest,
			wantCalls:  0,
		},
		{
			name:       "rating above range",
			body:       `{"place_id":"ChIJx","user_id":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","rating":5.01}`,
			wantStatus: http.StatusBadRequest,
			wantCalls:  0,
		},
		{
			name:       "rating below range",
			body:       `{"place_id":"ChIJx","user_id":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","rating":0.99}`,
			wantStatus: http.StatusBadRequest,
			wantCalls:  0,
		},
		{
			name:         "use case bad params",
			body:         validBody,
			stubErr:      domain.ErrBadParams,
			wantStatus:   http.StatusBadRequest,
			wantErrExact: domain.ErrBadParams.Error(),
			wantCalls:    1,
		},
		{
			name:         "use case not found",
			body:         validBody,
			stubErr:      domain.ErrNotFound,
			wantStatus:   http.StatusNotFound,
			wantErrExact: domain.ErrNotFound.Error(),
			wantCalls:    1,
		},
		{
			name:         "use case conflict",
			body:         validBody,
			stubErr:      domain.ErrConflict,
			wantStatus:   http.StatusConflict,
			wantErrExact: domain.ErrConflict.Error(),
			wantCalls:    1,
		},
		{
			name:         "use case unexpected",
			body:         validBody,
			stubErr:      errors.New("db unavailable"),
			wantStatus:   http.StatusInternalServerError,
			wantErrExact: domain.ErrUnexpected.Error(),
			wantCalls:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &stubReviewCreator{result: tt.stubResult, err: tt.stubErr}
			srv := httptest.NewServer(newReviewRouter(stub))
			defer srv.Close()

			req, err := http.NewRequest(http.MethodPost, srv.URL+reviewAPIPath, bytes.NewBufferString(tt.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

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

			if tt.wantTrimmedPlaceID != "" {
				assert.Equal(t, tt.wantTrimmedPlaceID, stub.lastIn.PlaceID)
			}
		})
	}
}
