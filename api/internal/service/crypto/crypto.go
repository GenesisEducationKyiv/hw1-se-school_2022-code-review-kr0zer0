package crypto

import (
	"api/internal/entities"
	"github.com/jellydator/ttlcache/v3"
	"time"
)

type (
	CachedCryptoService struct {
		cryptoService *CryptoService
		rateCache     *ttlcache.Cache[string, float64]
		cacheTTL      time.Duration
	}

	CryptoService struct {
		cryptoChain CryptoChain
	}

	Provider interface {
		GetExchangeRate(currencyPair entities.CurrencyPair) (float64, error)
	}

	ProviderCreator interface {
		CreateCryptoProvider() Provider
	}
)

func NewCachedCryptoService(cryptoService *CryptoService, cacheTTL time.Duration) *CachedCryptoService {
	rateCache := ttlcache.New[string, float64]()
	return &CachedCryptoService{
		cryptoService: cryptoService,
		rateCache:     rateCache,
		cacheTTL:      cacheTTL,
	}
}

func NewCryptoService(cryptoChain CryptoChain) *CryptoService {
	return &CryptoService{
		cryptoChain: cryptoChain,
	}
}

func (s *CryptoService) GetBtcUahRate() (float64, error) {
	pair := entities.NewCurrencyPair(entities.BTC, entities.UAH)
	return s.cryptoChain.HandleExchangeRate(pair)
}

func (c *CachedCryptoService) GetBtcUahRate() (float64, error) {
	cachedItem := c.rateCache.Get("rate", ttlcache.WithDisableTouchOnHit[string, float64]())
	if cachedItem != nil && !cachedItem.IsExpired() {
		return cachedItem.Value(), nil
	}

	rate, err := c.cryptoService.GetBtcUahRate()
	if err != nil {
		return -1, err
	}

	c.rateCache.Set("rate", rate, c.cacheTTL)

	return rate, nil
}