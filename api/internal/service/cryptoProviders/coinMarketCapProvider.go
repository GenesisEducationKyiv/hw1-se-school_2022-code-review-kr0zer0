package crypto_providers

import (
	"api/config"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/simonnilsson/ask"
)

type (
	CoinMarketCapProvider struct {
		HeaderName string
		APIKey     string
		APIUrl     string
	}

	CoinMarketCapProviderCreator struct {
		cfg *config.Config
	}
)

func NewCoinMarketCapProviderCreator(cfg *config.Config) *CoinMarketCapProviderCreator {
	return &CoinMarketCapProviderCreator{cfg: cfg}
}

func (c *CoinMarketCapProviderCreator) CreateCryptoProvider() *CoinMarketCapProvider {
	return &CoinMarketCapProvider{
		HeaderName: c.cfg.CryptoProviders.CoinMarketCap.HeaderName,
		APIKey:     c.cfg.CryptoProviders.CoinMarketCap.APIKey,
		APIUrl:     c.cfg.CryptoProviders.CoinMarketCap.URL,
	}
}

func (p *CoinMarketCapProvider) GetExchangeRate(baseCurrency, quoteCurrency string) (float64, error) {
	response, err := p.makeAPIRequest(baseCurrency, quoteCurrency)
	if err != nil {
		return -1, err
	}

	var mappedResponse map[string]interface{}
	err = json.Unmarshal(response, &mappedResponse)
	if err != nil {
		return -1, err
	}

	queryString := fmt.Sprintf("data.%s[0].quote.%s.price", baseCurrency, quoteCurrency)
	price, ok := ask.For(mappedResponse, queryString).Float(-1)
	if !ok {
		return price, fmt.Errorf("error when parsing JSON %v", mappedResponse)
	}

	return price, nil
}

func (p *CoinMarketCapProvider) makeAPIRequest(baseCurrency, quoteCurrency string) ([]byte, error) {
	client := resty.New()
	response, err := client.R().
		SetQueryParams(map[string]string{
			"symbol":  baseCurrency,
			"convert": quoteCurrency,
		}).
		SetHeader(p.HeaderName, p.APIKey).
		Get(p.APIUrl)
	if err != nil {
		return nil, err
	}

	return response.Body(), nil
}
