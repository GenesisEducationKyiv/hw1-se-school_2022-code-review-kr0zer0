package crypto_providers

import (
	"api/config"
	"api/internal/service"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type (
	CoinbaseProvider struct {
		APIUrl string
	}

	CoinbaseProviderCreator struct {
		cfg *config.Config
	}
)

func NewCoinbaseProviderCreator(cfg *config.Config) *CoinbaseProviderCreator {
	return &CoinbaseProviderCreator{cfg: cfg}
}

func (c *CoinbaseProviderCreator) CreateCryptoProvider() service.CryptoProvider {
	return &CoinbaseProvider{
		APIUrl: c.cfg.CryptoProviders.Coinbase.URL,
	}
}

type coinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}

func (p *CoinbaseProvider) GetExchangeRate(baseCurrency, quoteCurrency string) (float64, error) {
	response, err := p.makeAPIRequest(baseCurrency, quoteCurrency)
	if err != nil {
		return -1, err
	}

	var mappedResponse coinbaseResponse
	err = json.Unmarshal(response, &mappedResponse)
	if err != nil {
		return -1, err
	}

	price, err := strconv.ParseFloat(mappedResponse.Data.Amount, 64)
	if err != nil {
		return -1, fmt.Errorf("can't parse %v to float", price)
	}

	return price, nil
}

func (p *CoinbaseProvider) makeAPIRequest(baseCurrency, quoteCurrency string) ([]byte, error) {
	client := resty.New()
	response, err := client.R().
		SetPathParams(map[string]string{
			"base":  baseCurrency,
			"quote": quoteCurrency,
		}).
		Get(p.APIUrl)
	if err != nil {
		return nil, err
	}

	return response.Body(), nil
}
