package service

import (
	"api/config"
)

func NewService(repositories *Repository, cryptoChain CryptoChain, mailer Mailer, cfg *config.Config) *Service {
	crypto := NewCachedCryptoService(NewCryptoService(cryptoChain), cfg.Cache.RateCacheTTL)
	return &Service{
		Crypto: crypto,
		EmailSub: NewEmailSubscriptionService(
			repositories.EmailSubscriptionRepo,
			mailer,
			crypto,
		),
	}
}
