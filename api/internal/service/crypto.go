package service

type (
	CryptoService struct {
		cryptoProvider CryptoProvider
	}

	CryptoProvider interface {
		GetExchangeRate(baseCurrency, quoteCurrency string) (float64, error)
	}

	CryptoProviderCreator interface {
		CreateCryptoProvider() CryptoProvider
	}
)

func NewCryptoService(cryptoProvider CryptoProvider) *CryptoService {
	return &CryptoService{
		cryptoProvider: cryptoProvider,
	}
}

func (s *CryptoService) GetBtcUahRate() (float64, error) {
	return s.cryptoProvider.GetExchangeRate("BTC", "UAH")
}
