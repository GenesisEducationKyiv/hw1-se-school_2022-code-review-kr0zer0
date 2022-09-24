package service

import (
	"api/config"
	crypto2 "api/internal/service/crypto"
	"api/internal/service/emailSubscription"
	"api/internal/service/interfaces"
)

func NewService(repositories *interfaces.Repository, cryptoChain crypto2.CryptoChain, mailer interfaces.Mailer, cfg *config.Config) *interfaces.Service {
	crypto := crypto2.NewCachedCryptoService(crypto2.NewCryptoService(cryptoChain), cfg.Cache.RateCacheTTL)
	return &interfaces.Service{
		CryptoService: crypto,
		EmailSubService: emailSubscription.NewEmailSubscriptionService(
			repositories.EmailSubscriptionRepo,
			mailer,
			crypto,
		),
	}
}
