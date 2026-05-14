package app

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	v1 "post-con-back/internal/app/v1"
	"post-con-back/internal/domain"
	"post-con-back/internal/gateway/postgres/repositories"
)

func NewRouter(db *sql.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.String(http.StatusOK, "ok")
	})

	repo := repositories.NewReviewsRepository(db)
	stationRepo := repositories.NewStationRepository(db)
	uc := domain.NewReviewCreatorUseCase(repo, stationRepo)
	h := v1.NewReviewHandler(uc)
	apiv1 := r.Group("/api/v1")
	apiv1.POST("/review", h.CreateReview)

	return r
}
