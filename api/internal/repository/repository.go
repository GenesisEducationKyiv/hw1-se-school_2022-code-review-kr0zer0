package repository

import (
	"api/config"
	"api/internal/service"

	"github.com/mailjet/mailjet-apiv3-go"
)

func NewRepository(filepath string, cfg *config.Config, mailjetClient *mailjet.Client) *service.Repository {
	return &service.Repository{
		EmailSubscriptionRepo: NewEmailSubscriptionRepository(filepath),
		EmailSendingRepo:      NewEmailSendingRepository(cfg, mailjetClient),
	}
}
