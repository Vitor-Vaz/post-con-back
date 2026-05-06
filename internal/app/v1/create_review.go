package v1

import (
	"context"
	"errors"
	"net/http"
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
	PlaceID string  `json:"place_id"`
	UserID  string  `json:"user_id"`
	Rating  float64 `json:"rating"`
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	out, err := h.uc.CreateReview(c.Request.Context(), domain.CreateReviewInput{
		PlaceID: req.PlaceID,
		UserID:  uid,
		Rating:  req.Rating,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmptyPlaceID),
			errors.Is(err, domain.ErrInvalidUserID),
			errors.Is(err, domain.ErrInvalidReviewRating):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
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
