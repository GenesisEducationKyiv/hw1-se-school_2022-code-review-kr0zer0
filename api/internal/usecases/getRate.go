package usecases

import (
	"api/internal/entities"
	"api/internal/usecases/usecases_contracts"
)

type GetRateUseCase struct {
	cryptoProvider usecases_contracts.CryptoProvider
}

func NewGetRateUseCase(cryptoProvider usecases_contracts.CryptoProvider) *GetRateUseCase {
	return &GetRateUseCase{
		cryptoProvider: cryptoProvider,
	}
}

func (u *GetRateUseCase) GetBtcUahRate() (float64, error) {
	pair := entities.NewCurrencyPair(entities.BTC, entities.UAH)
	return u.cryptoProvider.GetExchangeRate(pair)
}
