package service

type ReviewSerive struct {
	repository IReviewRepository
}

type IReviewRepository interface{}

func NewUserService(repository IReviewRepository) *ReviewSerive {
	return &ReviewSerive{
		repository: repository,
	}
}
