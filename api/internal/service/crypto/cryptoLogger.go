package crypto

import (
	"api/internal/entities"
	"log"
	"reflect"
)

type LoggingCryptoProvider struct {
	cryptoProvider Provider
}

func NewLoggingCryptoProvider(provider Provider) *LoggingCryptoProvider {
	return &LoggingCryptoProvider{cryptoProvider: provider}
}

func (l *LoggingCryptoProvider) GetExchangeRate(currencyPair entities.CurrencyPair) (float64, error) {
	rate, err := l.cryptoProvider.GetExchangeRate(currencyPair)
	if err != nil {
		log.Printf("%v - %v", l.getProviderName(), err)
		return -1, err
	}

	log.Printf("%v - %v", l.getProviderName(), rate)
	return rate, nil
}

func (l *LoggingCryptoProvider) getProviderName() string {
	return reflect.TypeOf(l.cryptoProvider).Elem().Name()
}
