package service

type CryptoChain interface {
	SetNext(CryptoChain)
	HandleExchangeRate(baseCurrency, quoteCurrency string) (float64, error)
}

type BaseCryptoChain struct {
	next     CryptoChain
	provider CryptoProvider
}

func NewBaseCryptoChain(provider CryptoProvider) CryptoChain {
	return &BaseCryptoChain{provider: provider}
}

func (c *BaseCryptoChain) SetNext(next CryptoChain) {
	c.next = next
}

func (c *BaseCryptoChain) HandleExchangeRate(baseCurrency, quoteCurrency string) (float64, error) {
	rate, err := c.provider.GetExchangeRate(baseCurrency, quoteCurrency)
	if err != nil {
		nextChain := c.next
		if nextChain == nil {
			return -1, err
		}

		rate, err = nextChain.HandleExchangeRate(baseCurrency, quoteCurrency)
	}

	return rate, err
}
