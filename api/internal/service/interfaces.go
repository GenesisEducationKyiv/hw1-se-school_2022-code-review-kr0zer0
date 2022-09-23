package service

//go:generate mockgen -source=interfaces.go -destination=mocks/mocks.go

type (
	Crypto interface {
		GetBtcUahRate() (float64, error)
	}

	EmailSub interface {
		SendToAll() error
		Subscribe(email string) error
	}

	Service struct {
		Crypto
		EmailSub
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
