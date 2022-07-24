package repository

type EmailSubscription interface {
}

type Repository struct {
	EmailSubscription
}

func NewRepository() *Repository {
	return &Repository{
		EmailSubscription: NewEmailSubscriptionRepository(),
	}
}
