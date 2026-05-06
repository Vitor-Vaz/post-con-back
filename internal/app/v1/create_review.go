package v1

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"post-con-back/internal/domain"
)

type ReviewCreator interface {
	CreateReview(ctx context.Context, in domain.CreateReviewInput) (domain.Review, error)
}

type ReviewHandler struct {
	uc ReviewCreator
}

func NewReviewHandler(uc ReviewCreator) *ReviewHandler {
	return &ReviewHandler{uc: uc}
}

type createReviewRequest struct {
	PlaceID string    `json:"place_id" binding:"required"`
	UserID  uuid.UUID `json:"user_id" binding:"required"`
	Rating  float64   `json:"rating" binding:"required,gte=1,lte=5"`
}

type createReviewResponse struct {
	ID        string    `json:"id"`
	PlaceID   string    `json:"place_id"`
	UserID    string    `json:"user_id"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req createReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrBadRequest.Error()})
		return
	}
	placeID := strings.TrimSpace(req.PlaceID)
	if placeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrBadParams.Error()})
		return
	}
	out, err := h.uc.CreateReview(c.Request.Context(), domain.CreateReviewInput{
		PlaceID: placeID,
		UserID:  req.UserID,
		Rating:  req.Rating,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadParams):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrUnexpected.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, createReviewResponse{
		ID:        out.ID.String(),
		PlaceID:   out.PlaceID,
		UserID:    out.UserID.String(),
		Rating:    out.Rating,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	})
}
