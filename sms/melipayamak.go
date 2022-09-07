package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/entities"
	"net/http"
	"strconv"
)

type melipayamakProvider struct {
	entities.Provider
}

func makeRequest(jsonData map[string]string, op string) {
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("https://rest.payamak-panel.com/api/SendSMS/"+op, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func (p *melipayamakProvider) SendSMS(to, message string, isFlash bool) {
	jsonData := map[string]string{
		"username": p.Username,
		"password": p.Password,
		"to":       to,
		"from":     p.PhoneNumber,
		"text":     message,
		"isFlash":  strconv.FormatBool(isFlash),
	}
	go makeRequest(jsonData, "SendSMS")
}

func (p *melipayamakProvider) GetCredit() {
	jsonData := map[string]string{
		"UserName": p.Username,
		"PassWord": p.Password,
	}
	go makeRequest(jsonData, "GetCredit")
}
