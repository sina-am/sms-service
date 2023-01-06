package sender

type fakeSender struct{}

func NewfakeSender() *fakeSender {
	return &fakeSender{}
}

func (s *fakeSender) SendSMS(username, password, to, from, text string, isFlash bool) error {
	return nil
}

func (s *fakeSender) GetCredit(username, password string) error {
	return nil
}
func (s *fakeSender) Validate(username, password string) error {
	return nil
}
