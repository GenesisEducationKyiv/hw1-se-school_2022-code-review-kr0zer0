package service

import "api/internal/repository"

type Crypto interface {
	GetCurrentExchangeRate() (float64, error)
}

type Service struct {
	Crypto
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Crypto: NewCryptoService(),
	}
}
