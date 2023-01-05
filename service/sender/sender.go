package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/entities"
	"net/http"
	"strconv"
)

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

type SMSSender interface {
	SendSMS(username, password, to, from, text string, isFlash bool) error
	GetCredit(username, password string) error
	Validate(username, password string) error
}

type mockSender struct{}

func NewMockSender() *mockSender {
	return &mockSender{}
}

func (s *mockSender) SendSMS(username, password, to, from, text string, isFlash bool) error {
	return nil
}

func (s *mockSender) GetCredit(username, password string) error {
	return nil
}
func (s *mockSender) Validate(username, password string) error {
	return nil
}

type melipayamakSender struct {
	baseUrl string
}

func NewMelipayamakSender() *melipayamakSender {
	return &melipayamakSender{baseUrl: "https://rest.payamak-panel.com/api/SendSMS/"}
}

func (s *melipayamakSender) makeRequest(jsonData map[string]string, op string) (*http.Response, error) {
	jsonValue, _ := json.Marshal(jsonData)
	return http.Post(s.baseUrl+op, "application/json", bytes.NewBuffer(jsonValue))
}

func (s *melipayamakSender) SendSMS(username, password, to, from, text string, isFlash bool) error {
	jsonData := map[string]string{
		"username": username,
		"password": password,
		"to":       to,
		"from":     from,
		"text":     text,
		"isFlash":  strconv.FormatBool(isFlash),
	}
	_, err := s.makeRequest(jsonData, "SendSMS")
	return err
}

type getCreditResponse struct {
	Value        string
	ResStatus    int
	StrRetStatus string
}

func (s *melipayamakSender) GetCredit(username string, password string) error {
	jsonData := map[string]string{
		"UserName": username,
		"PassWord": password,
	}
	_, err := s.makeRequest(jsonData, "GetCredit")
	return err
}

func (s *melipayamakSender) Validate(username string, password string) error {
	jsonData := map[string]string{
		"UserName": username,
		"PassWord": password,
	}
	res, err := s.makeRequest(jsonData, "GetCredit")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	credit := &getCreditResponse{}
	if err := json.NewDecoder(res.Body).Decode(credit); err != nil {
		return err
	}

	if credit.StrRetStatus != "Ok" {
		return fmt.Errorf("%s", credit.StrRetStatus)
	}
	return nil
}

// func (s *melipayamakSender) GetDeliveries2(username string, password string, recID int64) {
// 	jsonData := map[string]string{
// 		"username": username,
// 		"password": password,
// 		"recID":    strconv.FormatInt(recID, 10),
// 	}
// 	go s.makeRequest(jsonData, "GetDeliveries2")
// }
// func (s *melipayamakSender) GetBasePrice(username string, password string) {
// 	jsonData := map[string]string{
// 		"UserName": username,
// 		"PassWord": password,
// 	}
// 	go s.makeRequest(jsonData, "GetBasePrice")
// }
