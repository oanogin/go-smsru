package smsru

import (
	"io"
	"net/http"
	"net/url"
)

const API_URL = "https://sms.ru"

type Client struct {
	APIID    string
	HTTP     *http.Client
	Test     bool
	JSON     bool
	Translit bool
}

func NewClient(aid string, testF, jsonF, translitF bool) *Client {
	return &Client{
		APIID:    aid,
		HTTP:     &http.Client{},
		Test:     testF,
		JSON:     jsonF,
		Translit: translitF,
	}
}

func (c *Client) makeRequest(endpoint string, params url.Values) (io.ReadCloser, error) {
	params.Set("api_id", c.APIID)
	url := API_URL + endpoint + "?" + params.Encode()

	resp, err := c.HTTP.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Body, nil
}
