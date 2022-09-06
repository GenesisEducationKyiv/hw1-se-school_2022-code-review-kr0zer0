package repository

import (
	"api/config"

	"github.com/mailjet/mailjet-apiv3-go"
)

type EmailSubscription interface {
	Add(email string) error
	GetAll() ([]string, error)
}

type EmailSending interface {
	SendToList(emails []string, message string) error
}

type Repository struct {
	EmailSubscription
	EmailSending
}

func NewRepository(filepath string, cfg *config.Config, mailjetClient *mailjet.Client) *Repository {
	return &Repository{
		EmailSubscription: NewEmailSubscriptionRepository(filepath),
		EmailSending:      NewEmailSendingRepository(cfg, mailjetClient),
	}
}
