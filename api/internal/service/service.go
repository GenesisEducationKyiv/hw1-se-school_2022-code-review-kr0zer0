package service

import "api/internal/repository"

type Crypto interface {
	GetCurrentExchangeRate() (float64, error)
}

type EmailSub interface {
	SendToAll() error
	Subscribe(email string) error
}

type Service struct {
	Crypto
	EmailSub
}

func NewService(repository *repository.Repository) *Service {
	crypto := NewCryptoService()
	return &Service{
		Crypto:   crypto,
		EmailSub: NewEmailSubscriptionService(repository.EmailSubscription, repository.EmailSending, crypto),
	}
}
