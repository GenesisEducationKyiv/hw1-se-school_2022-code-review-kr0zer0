package details

import (
	"api/internal/usecases"
	"github.com/jellydator/ttlcache/v3"
	"time"
)

type CachedRateGetter struct {
	getRateUseCase *usecases.GetRateUseCase
	rateCache      *ttlcache.Cache[string, float64]
	cacheTTL       time.Duration
}

func NewCachedRateGetter(cryptoService *usecases.GetRateUseCase, cacheTTL time.Duration) *CachedRateGetter {
	rateCache := ttlcache.New[string, float64]()
	return &CachedRateGetter{
		getRateUseCase: cryptoService,
		rateCache:      rateCache,
		cacheTTL:       cacheTTL,
	}
}

func (c *CachedRateGetter) GetBtcUahRate() (float64, error) {
	cachedItem := c.rateCache.Get("rate", ttlcache.WithDisableTouchOnHit[string, float64]())
	if cachedItem != nil && !cachedItem.IsExpired() {
		return cachedItem.Value(), nil
	}

	rate, err := c.getRateUseCase.GetBtcUahRate()
	if err != nil {
		return -1, err
	}

	c.rateCache.Set("rate", rate, c.cacheTTL)

	return rate, nil
}
