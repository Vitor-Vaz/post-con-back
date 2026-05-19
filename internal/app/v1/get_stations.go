package v1

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"post-con-back/internal/domain"
)

type StationsGetter interface {
	GetStations(ctx context.Context, page int) (domain.GetStationsOutput, error)
}

type StationsHandler struct {
	uc StationsGetter
}

func NewStationsHandler(uc StationsGetter) *StationsHandler {
	return &StationsHandler{uc: uc}
}

type stationResponse struct {
	ID          string   `json:"id"`
	PlaceID     string   `json:"place_id"`
	Name        string   `json:"name"`
	Address     *string  `json:"address,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
	TotalScore  float64  `json:"total_score"`
	ReviewCount int32    `json:"review_count"`
	Summary     *string  `json:"summary,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type paginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type getStationsResponse struct {
	Data       []stationResponse  `json:"data"`
	Pagination paginationResponse `json:"pagination"`
}

func (h *StationsHandler) GetStations(c *gin.Context) {
	page := 1
	if raw := c.Query("page"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrBadParams.Error()})
			return
		}
		page = parsed
	}

	out, err := h.uc.GetStations(c.Request.Context(), page)
	switch {
	case errors.Is(err, domain.ErrBadParams):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrUnexpected.Error()})
		return
	}

	data := make([]stationResponse, 0, len(out.Stations))
	for _, s := range out.Stations {
		data = append(data, stationResponse{
			ID:          s.ID.String(),
			PlaceID:     s.PlaceID,
			Name:        s.Name,
			Address:     s.Address,
			Latitude:    s.Latitude,
			Longitude:   s.Longitude,
			TotalScore:  s.TotalScore,
			ReviewCount: s.ReviewCount,
			Summary:     s.Summary,
			CreatedAt:   s.CreatedAt.Format(time.RFC3339Nano),
			UpdatedAt:   s.UpdatedAt.Format(time.RFC3339Nano),
		})
	}

	c.JSON(http.StatusOK, getStationsResponse{
		Data: data,
		Pagination: paginationResponse{
			Page:       out.Page,
			PageSize:   out.PageSize,
			Total:      out.Total,
			TotalPages: totalPages(out.Total, out.PageSize),
		},
	})
}

func totalPages(total int64, pageSize int) int {
	if total == 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}
