package service

import (
	"encoding/json"

	"github.com/kwiats/rate-all-things/pkg/model"
	"github.com/kwiats/rate-all-things/pkg/schema"
	"gorm.io/datatypes"
)

type ReviewService struct {
	Repository  IReviewRepository
	BlobStorage IBlobStorage
}

type IReviewRepository interface {
	CollectReviewItems() ([]model.ReviewItem, error)
	CollectReviewItem(uint) (model.ReviewItem, error)
	NewReviewItem(model.ReviewItem, []model.ReviewItemData) (bool, error)
}

type IBlobStorage interface {
	SaveFile([]byte) (string, error)
}

func NewReviewService(repository IReviewRepository, blobStorage IBlobStorage) *ReviewService {
	return &ReviewService{
		Repository:  repository,
		BlobStorage: blobStorage,
	}
}

func (service *ReviewService) ReviewItems() ([]*schema.ReviewItem, error) {
	reviewItems, err := service.Repository.CollectReviewItems()
	if err != nil {
		return nil, err
	}
	reviewItemsResult := make([]*schema.ReviewItem, 0, len(reviewItems))
	for _, item := range reviewItems {
		schemaItem := &schema.ReviewItem{
			ID:             item.ID,
			CategoryID:     item.CategoryID,
			UserID:         item.UserID,
			ReviewItemData: schema.MapReviewItemDataModelToDTO(item.ReviewItemData),
			Reviews:        schema.MapReviewsModelToDTO(item.Reviews),
		}
		reviewItemsResult = append(reviewItemsResult, schemaItem)
	}
	return reviewItemsResult, nil
}
func (service *ReviewService) ReviewItem(reviewId uint) (*schema.ReviewItem, error) {
	reviewItem, err := service.Repository.CollectReviewItem(reviewId)
	if err != nil {
		return nil, err
	}
	return &schema.ReviewItem{
		ID:             reviewItem.ID,
		CategoryID:     reviewItem.CategoryID,
		UserID:         reviewItem.UserID,
		ReviewItemData: schema.MapReviewItemDataModelToDTO(reviewItem.ReviewItemData),
		Reviews:        schema.MapReviewsModelToDTO(reviewItem.Reviews),
	}, nil
}

func (service *ReviewService) CreateReviewItem(userId uint, reviewItem schema.CreateReviewItem) (bool, error) {
	reviewItemDataDAO := make([]model.ReviewItemData, 0, len(reviewItem.ReviewItemData))
	for _, data := range reviewItem.ReviewItemData {
		var processedValue interface{} = data.Value
		switch v := data.Value.(type) {
		case []byte:
			filePath, err := service.BlobStorage.SaveFile(v)
			if err != nil {
				return false, err
			}
			processedValue = filePath
		default:
			processedValue = v
		}
		valueJSON, err := json.Marshal(processedValue)
		if err != nil {
			return false, err // Handle JSON marshaling error
		}
		reviewItemDataDAO = append(reviewItemDataDAO, model.ReviewItemData{
			CategoryCustomFieldID: data.CategoryCustomFieldID,
			Value:                 datatypes.JSON(valueJSON),
		})
	}

	reviewItemDAO := model.ReviewItem{
		UserID:     userId,
		CategoryID: reviewItem.CategoryID,
	}

	_, err := service.Repository.NewReviewItem(reviewItemDAO, reviewItemDataDAO)
	if err != nil {
		return false, err
	}
	return true, nil
}
