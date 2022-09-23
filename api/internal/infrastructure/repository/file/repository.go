package file

import (
	"api/internal/service"
)

func NewRepository(filepath string) *service.Repository {
	return &service.Repository{
		EmailSubscriptionRepo: NewEmailSubscriptionRepository(filepath),
	}
}
