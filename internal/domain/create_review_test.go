package domain_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"post-con-back/extension/testhelpers"
	"post-con-back/internal/domain"
	"post-con-back/internal/gateway/postgres/repositories"
)

func TestCreateReviewUseCase(t *testing.T) {
	tx := testhelpers.SetupTestDB(t)
	defer tx.Rollback()

	reviewsRepo := repositories.NewReviewsRepository(tx)
	stationRepo := repositories.NewStationRepository(tx)
	uc := domain.NewReviewCreatorUseCase(reviewsRepo, stationRepo)

	uid := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	t.Run("should create review and update station average on existing reviews", func(t *testing.T) {
		_, err := tx.Exec(`
			INSERT INTO reviews (place_id, user_id, rating) VALUES
			($1, $2, 4.0),
			($1, $3, 4.0),
			($1, $4, 4.0)
		`, "ChIJstation1", uuid.New(), uuid.New(), uuid.New())
		require.NoError(t, err)

		review, err := uc.CreateReview(context.Background(), domain.CreateReviewInput{
			PlaceID: "ChIJstation1",
			UserID:  uid,
			Rating:  3.0,
		})
		require.NoError(t, err)
		assert.Equal(t, "ChIJstation1", review.PlaceID)
		assert.Equal(t, 3.0, review.Rating)

		var totalScore float64
		var reviewCount int32
		var stationName string
		err = tx.QueryRow(`SELECT total_score, review_count, name FROM station WHERE place_id = $1`, "ChIJstation1").Scan(&totalScore, &reviewCount, &stationName)
		require.NoError(t, err)
		assert.Equal(t, 3.75, totalScore)
		assert.Equal(t, int32(4), reviewCount)
		assert.Equal(t, "ChIJstation1", stationName)
	})

	t.Run("should create review and create station on first review", func(t *testing.T) {
		review, err := uc.CreateReview(context.Background(), domain.CreateReviewInput{
			PlaceID: "ChIJstation2",
			UserID:  uid,
			Rating:  5.0,
		})
		require.NoError(t, err)
		assert.Equal(t, "ChIJstation2", review.PlaceID)
		assert.Equal(t, 5.0, review.Rating)

		var totalScore float64
		var reviewCount int32
		var stationName string
		err = tx.QueryRow(`SELECT total_score, review_count, name FROM station WHERE place_id = $1`, "ChIJstation2").Scan(&totalScore, &reviewCount, &stationName)
		require.NoError(t, err)
		assert.Equal(t, 5.0, totalScore)
		assert.Equal(t, int32(1), reviewCount)
		assert.Equal(t, "ChIJstation2", stationName)
	})

	t.Run("should not change station name on score update", func(t *testing.T) {
		_, err := tx.Exec(`
			INSERT INTO station (place_id, name, total_score, review_count)
			VALUES ($1, $2, 4.0, 1)
		`, "ChIJstation3", "Posto Original")
		require.NoError(t, err)
		_, err = tx.Exec(`INSERT INTO reviews (place_id, user_id, rating) VALUES ($1, $2, 4.0)`, "ChIJstation3", uuid.New())
		require.NoError(t, err)

		_, err = uc.CreateReview(context.Background(), domain.CreateReviewInput{
			PlaceID: "ChIJstation3",
			UserID:  uid,
			Rating:  5.0,
		})
		require.NoError(t, err)

		var name string
		err = tx.QueryRow(`SELECT name FROM station WHERE place_id = $1`, "ChIJstation3").Scan(&name)
		require.NoError(t, err)
		assert.Equal(t, "Posto Original", name)
	})
}
