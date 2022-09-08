package service

import (
	"api/config"
	"api/internal/handler"
)

type EmailSubscriptionRepo interface {
	Add(email string) error
	GetAll() ([]string, error)
}

type EmailSendingRepo interface {
	SendToList(emails []string, message string) error
}

type Repository struct {
	EmailSubscriptionRepo
	EmailSendingRepo
}

func NewService(repositories *Repository, cfg *config.Config) *handler.Service {
	crypto := NewCryptoService(cfg)
	return &handler.Service{
		CryptoService: crypto,
		EmailSubService: NewEmailSubscriptionService(
			repositories.EmailSubscriptionRepo,
			repositories.EmailSendingRepo, crypto,
		),
	}
}
