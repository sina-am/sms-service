package sender

import (
	"fmt"
	"main/entities"
)

type SMSSender interface {
	SendSMS(username, password, to, from, text string, isFlash bool) error
	GetCredit(username, password string) error
	Validate(username, password string) error
}

type SenderMap map[entities.ProviderType]SMSSender

func (s SenderMap) Validate(p *entities.Provider) error {
	if sender, found := s[p.Type]; found {
		return sender.Validate(p.Username, p.Password)
	}
	return fmt.Errorf("invalid provider")
}

func (s SenderMap) SendSMS(p *entities.Provider, text, to string, isFlash bool) error {
	if sender, found := s[p.Type]; found {
		return sender.SendSMS(p.Username, p.Password, to, p.PhoneNumber, text, isFlash)
	}
	return fmt.Errorf("invalid provider")
}
