package repository

import (
	"fmt"

	"github.com/kwiats/rate-all-things/pkg/model"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

// NewReviewItem implements service.IReviewRepository.
func (r *ReviewRepository) NewReviewItem(reviewItem model.ReviewItem, reviewItemData []model.ReviewItemData) (bool, error) {

	tx := r.db.Begin()

	if err := tx.Create(&reviewItem).Error; err != nil {
		fmt.Print(err)
		return false, err
	}
	fmt.Println(reviewItem.ID)

	for i := range reviewItemData {
		reviewItemData[i].ReviewItemID = reviewItem.ID
	}

	if err := tx.Create(&reviewItemData).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit().Error; err != nil {
		return false, err
	}

	return true, nil

}

func (r *ReviewRepository) CollectReviewItems() ([]model.ReviewItem, error) {
	var reviewItems []model.ReviewItem
	result := r.db.Preload("ReviewItemData.CategoryCustomField.CustomField").Find(&reviewItems)

	if result.Error != nil {
		return nil, result.Error
	}
	return reviewItems, nil
}

func (r *ReviewRepository) CollectReviewItem(reviewId uint) (model.ReviewItem, error) {
	var reviewItem model.ReviewItem
	result := r.db.Preload("ReviewItemData.CategoryCustomField.CustomField").Where("id = ?", reviewId).First(&reviewItem)

	if result.Error != nil {
		return reviewItem, result.Error
	}
	return reviewItem, nil
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}
