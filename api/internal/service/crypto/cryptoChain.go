package crypto

import "api/internal/entities"

type CryptoChain interface {
	SetNext(CryptoChain)
	HandleExchangeRate(currencyPair entities.CurrencyPair) (float64, error)
}

type BaseCryptoChain struct {
	next     CryptoChain
	provider Provider
}

func NewBaseCryptoChain(provider Provider) CryptoChain {
	return &BaseCryptoChain{provider: provider}
}

func (c *BaseCryptoChain) SetNext(next CryptoChain) {
	c.next = next
}

func (c *BaseCryptoChain) HandleExchangeRate(currencyPair entities.CurrencyPair) (float64, error) {
	rate, err := c.provider.GetExchangeRate(currencyPair)
	if err != nil {
		nextChain := c.next
		if nextChain == nil {
			return -1, err
		}

		rate, err = nextChain.HandleExchangeRate(currencyPair)
	}

	return rate, err
}
