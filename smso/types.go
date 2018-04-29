package sms

import (
	"net/http"
	"time"
)

// Basic
type Client struct {
	ApiId string       `json:"api_id"`
	Http  *http.Client `json:"-"`
	Debug bool         `json:"-"`
}

type Response struct {
	Status     string            `json:"status"`
	StatusCode int               `json:"status_code"`
	Ids        []string          `json:"id"`
	Cost       float32           `json:"cost"`
	Count      int               `json:"count"`
	Balance    float32           `json:"balance"`
	Limit      int               `json:"limit"`
	LimitSent  int               `json:"limit_sent"`
	Senders    []string          `json:"senders"`
	Stoplist   map[string]string `json:"stoplist"`
	Callbacks  []string          `json:"callbacks"`
}

// Request describe sms
type Request struct {
	To        string            `json:"to"`
	Msg       string            `json:"msg"`
	Multi     map[string]string `json:"multi"`
	JSON      bool              `json:"json"`
	From      string            `json:"from"`
	Time      time.Time         `json:"time"`
	Translit  bool              `json:"translit"`
	Test      bool              `json:"test"`
	PartnerID int               `json:"partner_id"`
}
