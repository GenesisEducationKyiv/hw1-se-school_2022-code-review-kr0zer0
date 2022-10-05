package crypto_providers

import (
	"api/internal/entities"
	"api/internal/usecases/usecases_contracts"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

type LoggingCryptoProvider struct {
	cryptoProvider usecases_contracts.CryptoProvider
	logger         *logrus.Logger
}

func NewLoggingCryptoProvider(cryptoProvider usecases_contracts.CryptoProvider, logger *logrus.Logger) *LoggingCryptoProvider {
	return &LoggingCryptoProvider{
		cryptoProvider: cryptoProvider,
		logger:         logger,
	}
}

func (l *LoggingCryptoProvider) GetExchangeRate(currencyPair entities.CurrencyPair) (*entities.Rate, error) {
	rate, err := l.cryptoProvider.GetExchangeRate(currencyPair)
	if err != nil {
		l.logger.Info(fmt.Sprintf("%v - %v", l.getProviderName(), err))
		return nil, err
	}

	l.logger.Info(fmt.Sprintf("%v - %v", l.getProviderName(), rate.String()))
	return rate, nil
}

func (l *LoggingCryptoProvider) getProviderName() string {
	return reflect.TypeOf(l.cryptoProvider).Elem().Name()
}
