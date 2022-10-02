package crypto_providers

import (
	"api/config"
	"api/internal/constants"
	"api/internal/entities"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
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

func (c *CoinbaseProviderCreator) CreateCryptoProvider() *CoinbaseProvider {
	return &CoinbaseProvider{
		APIUrl: c.cfg.CryptoProviders.Coinbase.URL,
	}
}

type coinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}

func (p *CoinbaseProvider) GetExchangeRate(currencyPair entities.CurrencyPair) (*entities.Rate, error) {
	response, err := p.makeAPIRequest(string(currencyPair.GetBase()), string(currencyPair.GetQuote()))
	if err != nil {
		return nil, err
	}

	var mappedResponse coinbaseResponse
	err = json.Unmarshal(response, &mappedResponse)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(mappedResponse.Data.Amount, constants.Float64Size)
	if err != nil {
		return nil, fmt.Errorf("can't parse %v to float", price)
	}

	return entities.NewRate(currencyPair, price), nil
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
