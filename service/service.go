package service

import (
	"fmt"
	"main/database"
	"main/entities"
	"main/service/sender"
)

type SMSService interface {
	SendSMSWith(providerId int, message string, phoneNumber string) error
	SendSMS(message string, phoneNumber string) error
	CreateProvider(*entities.Provider) error
	UpdateProvider(*entities.Provider) error
}

type smsService struct {
	storage database.Storage
	senders sender.SenderMap
}

func NewSMSService(storage database.Storage, m sender.SenderMap) *smsService {
	return &smsService{storage: storage, senders: m}
}

func (s *smsService) CreateProvider(p *entities.Provider) error {
	// Check if already exist
	found, err := s.storage.ExistProviderByUsername(p.Username)
	if err != nil {
		return err
	}
	if found {
		return fmt.Errorf("provider already exist")
	}
	// Check if credentials are valid
	if err := s.senders.Validate(p); err != nil {
		return err
	}
	return s.storage.CreateProvider(p)
}

func (s *smsService) UpdateProvider(p *entities.Provider) error {
	return s.storage.UpdateProvider(p)
}

func (s *smsService) SendSMSWith(providerId int, message string, phoneNumber string) error {
	provider, err := s.storage.GetProviderById(providerId)
	if err != nil {
		return err
	}
	return s.sendSMS(provider, message, phoneNumber)
}

func (s *smsService) SendSMS(message string, phoneNumber string) error {
	providers, err := s.storage.GetAllProviders()
	if err != nil {
		return err
	}

	for _, provider := range providers {
		if err := s.sendSMS(provider, message, phoneNumber); err == nil {
			break
		}
	}
	return nil
}

func (s *smsService) sendSMS(p *entities.Provider, message string, phoneNumber string) error {
	err := s.senders.SendSMS(p, message, phoneNumber, false)
	switch err {
	case sender.ErrInvalidCredentials:
		p.InvalidCredential = true
		return s.UpdateProvider(p)
	case sender.ErrInsufficientCredit:
		p.InSufficientCredit = true
		return s.UpdateProvider(p)
	case sender.ErrProviderProblem:
		p.Failed += 1
		return s.UpdateProvider(p)
	default:
		p.Success += 1
		return s.UpdateProvider(p)
	}
}
