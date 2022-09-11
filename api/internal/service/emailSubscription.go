package service

import (
	"fmt"
)

type EmailSubscriptionService struct {
	emailSubsRepo    EmailSubscriptionRepo
	emailSendingRepo EmailSendingRepo
	cryptoService    *CryptoService
}

func NewEmailSubscriptionService(
	emailSubsRepo EmailSubscriptionRepo,
	emailSendingRepo EmailSendingRepo,
	cryptoService *CryptoService) *EmailSubscriptionService {
	return &EmailSubscriptionService{
		emailSubsRepo:    emailSubsRepo,
		emailSendingRepo: emailSendingRepo,
		cryptoService:    cryptoService,
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

	err = s.emailSendingRepo.SendToList(emails, fmt.Sprintf("%.2fUAH", rate))
	if err != nil {
		return err
	}

	return nil
}
