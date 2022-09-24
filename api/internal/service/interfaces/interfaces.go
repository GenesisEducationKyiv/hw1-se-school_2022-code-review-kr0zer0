package interfaces

//go:generate mockgen -source=interfaces.go -destination=mocks/mocks.go

type (
	CryptoService interface {
		GetBtcUahRate() (float64, error)
	}

	EmailSubService interface {
		SendToAll() error
		Subscribe(email string) error
	}

	Service struct {
		CryptoService
		EmailSubService
	}
)

type (
	EmailSubscriptionRepo interface {
		Add(email string) error
		GetAll() ([]string, error)
	}

	Repository struct {
		EmailSubscriptionRepo
	}

	Mailer interface {
		SendToList(emails []string, message string) error
	}
)
