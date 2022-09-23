package service

import (
	"fmt"
)

type EmailSubscriptionService struct {
	emailSubsRepo EmailSubscriptionRepo
	mailer        Mailer
	cryptoService Crypto
}

func NewEmailSubscriptionService(
	emailSubsRepo EmailSubscriptionRepo,
	mailer Mailer,
	cryptoService Crypto) *EmailSubscriptionService {
	return &EmailSubscriptionService{
		emailSubsRepo: emailSubsRepo,
		mailer:        mailer,
		cryptoService: cryptoService,
	}
}

func (s *EmailSubscriptionService) Subscribe(email string) error {
	err := s.emailSubsRepo.Add(email)
	if err != nil {
		return err
	}

	return nil
}

func (s *EmailSubscriptionService) SendToAll() error {
	emails, err := s.emailSubsRepo.GetAll()
	if err != nil {
		return err
	}

	rate, err := s.cryptoService.GetBtcUahRate()
	if err != nil {
		return err
	}

	err = s.mailer.SendToList(emails, fmt.Sprintf("%.2fUAH", rate))
	if err != nil {
		return err
	}

	return nil
}
