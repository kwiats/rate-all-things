package schema

import (
	"time"

	"github.com/kwiats/rate-all-things/pkg/model"
	"gorm.io/datatypes"
)

type ReviewDTO struct {
	ID            uint           `json:"id"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	Rate          int            `json:"rate"`
	CategoryID    uint           `json:"categoryId"`
	UserID        uint           `json:"userId"`
	CustomFieldID uint           `json:"customFieldId"`
	ReviewData    datatypes.JSON `json:"reviewData"`
}

func NewReviewDTO(review model.Review, fieldData model.ReviewItemData) ReviewDTO {
	return ReviewDTO{
		ID:        review.ID,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}
