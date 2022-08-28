package service

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/simonnilsson/ask"
)

type CryptoService struct {
}

func NewCryptoService() *CryptoService {
	return &CryptoService{}
}

func (s *CryptoService) GetCurrentExchangeRate() (float64, error) {
	url := "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest"
	client := resty.New()
	response, err := client.R().
		SetQueryParams(map[string]string{
			"symbol":  "BTC",
			"convert": "UAH",
		}).
		SetHeader("X-CMC_PRO_API_KEY", os.Getenv("COINMARKETCAP_API_KEY")).
		Get(url)
	if err != nil {
		return 0, err
	}

	var mappedResponse map[string]interface{}
	err = json.Unmarshal(response.Body(), &mappedResponse)
	if err != nil {
		return 0, err
	}

	price, ok := ask.For(mappedResponse, "data.BTC[0].quote.UAH.price").Float(0)
	if !ok {
		return price, errors.New("incorrect path when parsing JSON")
	}

	return price, nil
}
