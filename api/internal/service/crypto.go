package service

import (
	"api/config"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/simonnilsson/ask"
)

type CryptoService struct {
	cryptoProvider CryptoProvider
}

func NewCryptoService(cryptoProvider CryptoProvider) *CryptoService {
	return &CryptoService{
		cryptoProvider: cryptoProvider,
	}
}

func (s *CryptoService) GetBtcUahRate() (float64, error) {
	return s.cryptoProvider.GetExchangeRate("BTC", "UAH")
}

type CryptoProvider interface {
	GetExchangeRate(baseCurrency, quoteCurrency string) (float64, error)
}

type CoinMarketCapProvider struct {
	cfg *config.Config
}

func NewCoinMarketCapProvider(cfg *config.Config) *CoinMarketCapProvider {
	return &CoinMarketCapProvider{cfg: cfg}
}

func (p *CoinMarketCapProvider) makeAPIRequest(baseCurrency, quoteCurrency string) ([]byte, error) {
	client := resty.New()
	response, err := client.R().
		SetQueryParams(map[string]string{
			"symbol":  baseCurrency,
			"convert": quoteCurrency,
		}).
		SetHeader(p.cfg.CryptoAPI.HeaderName, p.cfg.CryptoAPI.APIKey).
		Get(p.cfg.CryptoAPI.URL)
	if err != nil {
		return nil, err
	}

	return response.Body(), nil
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
		return price, errors.New("incorrect path when parsing JSON")
	}

	return price, nil
}
