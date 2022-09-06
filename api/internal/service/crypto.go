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
	cfg *config.Config
}

func NewCryptoService(cfg *config.Config) *CryptoService {
	return &CryptoService{
		cfg: cfg,
	}
}

func (s *CryptoService) GetCurrentExchangeRate(cryptoSymbol, fiatSymbol string) (float64, error) {
	client := resty.New()
	response, err := client.R().
		SetQueryParams(map[string]string{
			"symbol":  cryptoSymbol,
			"convert": fiatSymbol,
		}).
		SetHeader(s.cfg.CryptoAPI.HeaderName, s.cfg.CryptoAPI.APIKey).
		Get(s.cfg.CryptoAPI.URL)
	if err != nil {
		return 0, err
	}

	var mappedResponse map[string]interface{}
	err = json.Unmarshal(response.Body(), &mappedResponse)
	if err != nil {
		return 0, err
	}

	queryString := fmt.Sprintf("data.%s[0].quote.%s.price", cryptoSymbol, fiatSymbol)
	price, ok := ask.For(mappedResponse, queryString).Float(0)
	if !ok {
		return price, errors.New("incorrect path when parsing JSON")
	}

	return price, nil
}

func (s *CryptoService) GetBtcUahRate() (float64, error) {
	return s.GetCurrentExchangeRate("BTC", "UAH")
}
