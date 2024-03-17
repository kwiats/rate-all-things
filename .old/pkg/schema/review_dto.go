package schema

import (
	"github.com/kwiats/rate-all-things/pkg/model"
	"gorm.io/datatypes"
)

type Review struct {
	ID           uint   `json:"id"`
	ReviewItemID uint   `json:"review_item_id"`
	UserID       uint   `json:"user_id"`
	Rating       int    `json:"rating"`
	Comment      string `json:"comment"`
}

type CreateReview struct {
	ReviewItemID uint   `json:"review_item_id"`
	UserID       uint   `json:"user_id"`
	Rating       int    `json:"rating"`
	Comment      string `json:"comment"`
}
type UpdateReview struct {
	Rating  int    `json:"rating,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type ReviewItemData struct {
	ID                  uint                   `json:"id"`
	ReviewItemID        uint                   `json:"review_item_id"`
	CategoryCustomField CategoryCustomFieldDTO `json:"category_custom_field"`
	Value               datatypes.JSON         `json:"value"`
}

type CreateReviewItemData struct {
	CategoryCustomFieldID uint        `json:"category_custom_field_id"`
	Value                 interface{} `json:"value"`
}
type UpdateReviewItemData struct {
	Value datatypes.JSON `json:"value,omitempty"`
}

type ReviewItem struct {
	ID             uint             `json:"id"`
	CategoryID     uint             `json:"category_id"`
	UserID         uint             `json:"user_id"`
	ReviewItemData []ReviewItemData `json:"review_item_data,omitempty"`
	Reviews        []Review         `json:"reviews,omitempty"`
}
type CreateReviewItem struct {
	CategoryID     uint                   `json:"category_id"`
	UserID         uint                   `json:"user_id"`
	ReviewItemData []CreateReviewItemData `json:"review_item_data,omitempty"`
}
type UpdateReviewItem struct {
	CategoryID     uint                   `json:"category_id,omitempty"`
	ReviewItemData []UpdateReviewItemData `json:"review_item_data,omitempty"`
}

func MapReviewItemDataModelToDTO(data []model.ReviewItemData) []ReviewItemData {
	dtoData := make([]ReviewItemData, len(data))
	for i, d := range data {
		dtoData[i] = ReviewItemData{
			ID:           d.ID,
			ReviewItemID: d.ReviewItemID,
			CategoryCustomField: CategoryCustomFieldDTO{
				ID:    d.CategoryCustomFieldID,
				Type:  d.CategoryCustomField.CustomField.Type,
				Title: d.CategoryCustomField.Title,
			},
			Value: d.Value,
		}
	}
	return dtoData
}
func MapReviewsModelToDTO(reviews []model.Review) []Review {
	dtoReviews := make([]Review, len(reviews))
	for i, r := range reviews {
		dtoReviews[i] = Review{
			ID:           r.ID,
			ReviewItemID: r.ReviewItemID,
			UserID:       r.UserID,
			Rating:       r.Rating,
			Comment:      r.Comment,
		}
	}
	return dtoReviews
}
