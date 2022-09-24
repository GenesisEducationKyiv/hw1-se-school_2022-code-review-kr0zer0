package file

import (
	"api/internal/service/interfaces"
)

func NewRepository(filepath string) *interfaces.Repository {
	return &interfaces.Repository{
		EmailSubscriptionRepo: NewEmailSubscriptionRepository(filepath),
	}
}
