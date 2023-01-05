package server

import (
	"net/http"
)

type Authenticator interface {
	Authenticate(r *http.Request) error
}

type noneAuthenticator struct {
}

func NewNoneAuthenticator() *noneAuthenticator {
	return &noneAuthenticator{}
}

func (auth *noneAuthenticator) Authenticate(r *http.Request) error {
	/* Always authenticate*/
	return nil
}
