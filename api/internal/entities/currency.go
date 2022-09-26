package entities

import "fmt"

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

func (p *CurrencyPair) String() string {
	return fmt.Sprintf("%v%v", p.Quote, p.Base)
}
