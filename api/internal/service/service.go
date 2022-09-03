package service

import (
	"api/config"
	"api/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Crypto interface {
	GetCurrentExchangeRate(cryptoSymbol, fiatSymbol string) (float64, error)
	GetBtcUahRate() (float64, error)
}

type EmailSub interface {
	SendToAll() error
	Subscribe(email string) error
}

type Service struct {
	Crypto
	EmailSub
}

func NewService(repositories *repository.Repository, cfg *config.Config) *Service {
	crypto := NewCryptoService(cfg)
	return &Service{
		Crypto:   crypto,
		EmailSub: NewEmailSubscriptionService(repositories.EmailSubscription, repositories.EmailSending, crypto),
	}
}
