package smsru

import "time"

// SendRequest example
/*
	https://sms.ru/sms/send?api_id=[APIID]&to=79255070602,74993221627&msg=hello+world&json=1
	https://sms.ru/sms/send?api_id=[APIID]&to[79255070602]=hello+world&to[74993221627]=hello+world&json=1
*/
type SendRequest struct {
	To        string
	Msg       string
	JSON      bool
	From      string
	Time      time.Time
	Translit  bool
	Test      bool
	PartnerID int
}

// SendResponse example
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
		"balance": 1000.76
	}
*/
type SendResponse struct {
	Status     string                   `json:"status"`
	StatusCode int                      `json:"status_code"`
	Sms        map[string]SendSmsResult `json:"sms"`
	Balance    float64                  `json:"balance"`
}

// SendSmsResult contain result about sms messages
type SendSmsResult struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SmsID      string `json:"sms_id"`
	StatusText string `json:"status_text"`
}
