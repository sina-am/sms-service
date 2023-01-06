package entities

import (
	"github.com/go-playground/validator/v10"
)

type ProviderType string

var validate *validator.Validate

func NewValidator() {
	validate = validator.New()
}

const (
	Melipayamak ProviderType = "Melipayamak"
)

type Provider struct {
	Id                 int          `json:"id"`
	Type               ProviderType `json:"type" `
	Username           string       `json:"username"`
	Password           string       `json:"password"`
	PhoneNumber        string       `json:"phone_number"`
	Success            int          `json:"success"`
	Failed             int          `json:"failed"`
	InvalidCredential  bool         `json:"invalid_credntial"`
	InSufficientCredit bool         `json:"insufficient_credit"`
}

type ProviderCreationRequest struct {
	Type        ProviderType `json:"type" validate:"required"`
	Username    string       `json:"username" validate:"required"`
	Password    string       `json:"password" validate:"required"`
	PhoneNumber string       `json:"phone_number" validate:"required"`
}

func (p *ProviderCreationRequest) Validate() error {
	return validate.Struct(p)
}

func (p *ProviderCreationRequest) ToSchema() *Provider {
	return &Provider{
		Type:        p.Type,
		Username:    p.Username,
		Password:    p.Password,
		PhoneNumber: p.PhoneNumber,
	}
}

type SendSMSRequest struct {
	Message     string `json:"message" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

func (r *SendSMSRequest) Validate() error {
	return validate.Struct(r)
}
