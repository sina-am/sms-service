package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type melipayamakSender struct {
	baseUrl string
}

func NewMelipayamakSender() *melipayamakSender {
	return &melipayamakSender{baseUrl: "https://rest.payamak-panel.com/api/SendSMS/"}
}

type melipayamakResponse struct {
	Value        string
	ResStatus    int
	StrRetStatus string
}

func (s *melipayamakSender) makeRequest(jsonData map[string]string, op string) (*melipayamakResponse, error) {
	jsonValue, _ := json.Marshal(jsonData)
	res, err := http.Post(s.baseUrl+op, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBody := &melipayamakResponse{}
	if err := json.NewDecoder(res.Body).Decode(resBody); err != nil {
		return nil, err
	}
	return resBody, nil
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
	res, err := s.makeRequest(jsonData, "SendSMS")
	if err != nil {
		return err
	}
	switch res.ResStatus {
	case 0:
		return ErrInvalidCredentials
	case 1:
		return s.GetDeliveries2(username, password, res.Value)
	case 2:
		return ErrInsufficientCredit
	default:
		return ErrProviderProblem
	}
}
func (s *melipayamakSender) GetDeliveries2(username string, password string, recID string) error {
	jsonData := map[string]string{
		"username": username,
		"password": password,
		"recID":    recID,
	}
	res, err := s.makeRequest(jsonData, "GetDeliveries2")
	if err != nil {
		return err
	}
	switch res.ResStatus {
	case 0:
		return nil
	case 1:
		return ErrProviderProblem
	case 3:
		return ErrInvalidCredentials
	default:
		return nil
	}
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
	credit, err := s.makeRequest(jsonData, "GetCredit")
	if err != nil {
		return err
	}
	if credit.StrRetStatus != "Ok" {
		return fmt.Errorf("%s", credit.StrRetStatus)
	}
	return nil
}
