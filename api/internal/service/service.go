package service

import (
	"api/config"
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

func NewService(repositories *Repository, cryptoChain CryptoChain, cfg *config.Config) *handler.Service {
	crypto := NewCachedCryptoService(NewCryptoService(cryptoChain), cfg.Cache.RateCacheTTL)
	return &handler.Service{
		CryptoService: crypto,
		EmailSubService: NewEmailSubscriptionService(
			repositories.EmailSubscriptionRepo,
			repositories.EmailSendingRepo, crypto,
		),
	}
}
