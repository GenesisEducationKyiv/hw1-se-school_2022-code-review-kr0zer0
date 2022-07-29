package repository

type EmailSubscription interface {
	Add(email string) error
	GetAll() ([]string, error)
	CheckIfExists(emailToFind string) (bool, error)
}

type EmailSending interface {
	SendToList(emails []string, message string) error
}

type Repository struct {
	EmailSubscription
	EmailSending
}

func NewRepository(filepath string) *Repository {
	return &Repository{
		EmailSubscription: NewEmailSubscriptionRepository(filepath),
		EmailSending:      NewEmailSendingRepository(),
	}
}
