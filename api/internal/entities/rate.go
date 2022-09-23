package entities

type Currency string

const (
	BTC Currency = "BTC"
	UAH Currency = "UAH"
)

type CurrencyPair struct {
	Base  Currency
	Quote Currency
}

func NewCurrencyPair(base, quote Currency) CurrencyPair {
	return CurrencyPair{
		Base:  base,
		Quote: quote,
	}
}
