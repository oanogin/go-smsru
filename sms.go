package smsru

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// SendSmsRequest example
/*
	https://sms.ru/sms/send?api_id=[APIID]&to=79255070602,74993221627&msg=hello+world&json=1
	https://sms.ru/sms/send?api_id=[APIID]&to[79255070602]=hello+world&to[74993221627]=hello+world&json=1
*/

// SendSmsResponse example
/*
	{
		"status": "OK",
		"status_code": 100,
		"sms": {
			"79999999999": {
				"status": "OK",
				"status_code": 100,
				"sms_id": "200000-1000000"
			},
			"74959999999": {
				"status": "ERROR",
				"status_code": 202,
				"status_text": "Неправильно указан номер телефона получателя, либо на него нет маршрута"
			}
		},
		"balance": 1000.74
	}
*/
type SendSmsResponse struct {
	Status     string                   `json:"status"`
	StatusCode int                      `json:"status_code"`
	Sms        map[string]SendSmsResult `json:"sms"`
	Balance    float64                  `json:"balance"`
}

// SendSmsResult contain result about sms messages inside SendResponse
type SendSmsResult struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SmsID      string `json:"sms_id"`
	StatusText string `json:"status_text"`
}

func (c *Client) SendSms(to, msg string) (*SendSmsResponse, error) {
	urlValues := url.Values{}
	urlValues.Set("to", HairPhone(to))
	urlValues.Set("msg", msg)

	if c.Test {
		urlValues.Set("test", "1")
	}

	if c.JSON {
		urlValues.Set("json", "1")
	}

	if c.Translit {
		urlValues.Set("translit", "1")
	}

	ssr := &SendSmsResponse{}

	result, err := c.makeRequest("/sms/send", urlValues)
	if err != nil {
		return nil, fmt.Errorf("MakeRequest func error %v", err)
	}

	err = json.Unmarshal(result.Bytes(), &ssr)
	if err != nil {
		return nil, fmt.Errorf("Unmarshalling error %v", err)
	}

	if ssr.Status != "OK" {
		return nil, fmt.Errorf("Sms.ru Error %d - %s", ssr.StatusCode, CODES[ssr.StatusCode])
	}

	return ssr, nil
}
