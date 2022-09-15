package crypto_providers

import (
	"api/internal/service"
	"log"
	"reflect"
)

type LoggingCryptoProvider struct {
	cryptoProvider service.CryptoProvider
}

func NewLoggingCryptoProvider(provider service.CryptoProvider) *LoggingCryptoProvider {
	return &LoggingCryptoProvider{cryptoProvider: provider}
}

func (l *LoggingCryptoProvider) GetExchangeRate(baseCurrency, quoteCurrency string) (float64, error) {
	rate, err := l.cryptoProvider.GetExchangeRate(baseCurrency, quoteCurrency)
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
