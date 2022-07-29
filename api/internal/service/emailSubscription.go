package service

import (
	"api/internal/repository"
	"errors"
	"fmt"
)

type EmailSubscriptionService struct {
	emailSubsRepo    repository.EmailSubscription
	emailSendingRepo repository.EmailSending
	cryptoService    Crypto
}

func NewEmailSubscriptionService(emailSubsRepo repository.EmailSubscription, emailSendingRepo repository.EmailSending, cryptoService Crypto) *EmailSubscriptionService {
	return &EmailSubscriptionService{
		emailSubsRepo:    emailSubsRepo,
		emailSendingRepo: emailSendingRepo,
		cryptoService:    cryptoService,
	}
}

var EmailDuplErr = errors.New("this email already exists")

func (s *EmailSubscriptionService) Subscribe(email string) error {
	exists, err := s.emailSubsRepo.CheckIfExists(email)
	if err != nil {
		return err
	}

	if exists {
		return EmailDuplErr
	}

	err = s.emailSubsRepo.Add(email)
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

	rate, err := s.cryptoService.GetCurrentExchangeRate()
	if err != nil {
		return err
	}

	err = s.emailSendingRepo.SendToList(emails, fmt.Sprintf("%.2f", rate))
	if err != nil {
		return err
	}

	return nil
}
