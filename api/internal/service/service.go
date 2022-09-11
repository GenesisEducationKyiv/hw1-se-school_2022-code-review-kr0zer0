package service

import (
	"api/internal/handler"
)

//go:generate mockgen -source=service.go -destination=mocks/repoMock.go

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

func NewService(repositories *Repository, cryptoProvider CryptoProvider) *handler.Service {
	crypto := NewCryptoService(cryptoProvider)
	return &handler.Service{
		CryptoService: crypto,
		EmailSubService: NewEmailSubscriptionService(
			repositories.EmailSubscriptionRepo,
			repositories.EmailSendingRepo, crypto,
		),
	}
}
