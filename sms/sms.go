package sms

import "main/entities"

type SMSProvider interface {
	SendSMS(to, message string, isFlash bool)
	GetCredit()
}

func Factory(provider entities.Provider) SMSProvider {
	switch provider.Type {
	case entities.Melipayamak:
		return &melipayamakProvider{provider}
	}
	return nil
}
