package sender

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInsufficientCredit = errors.New("insufficient credit")
	ErrRetry              = errors.New("retry another time")
	ErrProviderProblem    = errors.New("insufficient resource")
)
