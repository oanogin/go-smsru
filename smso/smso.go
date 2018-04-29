package smso

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

// Base API url

// NewClient creates a new Client instance.
//
// id is your api_id
func NewClient(id string) *Client {
	return NewClientWithHttp(id, &http.Client{})
}

// NewClientWithHttp creates a new Client instance
//
// and allows you to pass a http.Client.
func NewClientWithHttp(id string, client *http.Client) *Client {
	c := &Client{
		ApiId: id,
		Http:  client,
	}

	return c
}

// NewSms creates a new message
//
// to is where to send it (phone number), msg is the message text.
func NewSms(to string, msg string) *Request {
	return &Request{
		To:  to,
		Msg: msg,
	}
}

// NewMulti creates a one request for multiple messages
func NewMulti(sms ...*Request) *Request {
	arr := make(map[string]string)
	for _, o := range sms {
		arr[o.To] = o.Msg
	}

	return &Request{
		Multi: arr,
	}
}

func (c *Client) makeRequest(endpoint string, params url.Values) (Response, []string, error) {
	params.Set("api_id", c.ApiId)
	url := API_URL + endpoint + "?" + params.Encode()

	resp, err := c.Http.Get(url)
	if err != nil {
		return Response{}, nil, err
	}
	defer resp.Body.Close()

	sc := bufio.NewScanner(resp.Body)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return Response{}, nil, error_internal
	}

	if len(lines) == 0 {
		return Response{}, nil, error_no_response
	}

	status, _ := strconv.Atoi(lines[0])

	if status >= 200 {
		msg := fmt.Sprintf("Code: %d; Status: %s", status, codeStatus[status])
		return Response{}, nil, errors.New(msg)
	}

	res := Response{StatusCode: status}
	return res, lines, nil
}

// SmsSend will send a Sms item to Service
func (c *Client) SmsSend(p *Request) (Response, error) {
	var params = url.Values{}

	if len(p.Multi) > 0 {
		for to, text := range p.Multi {
			key := fmt.Sprintf("multi[%s]", to)
			params.Add(key, text)
		}
	} else {
		params.Set("to", p.To)
		params.Set("msg", p.Msg)
	}

	if len(p.From) > 0 {
		params.Set("from", p.From)
	}

	if p.PartnerID > 0 {
		val := strconv.Itoa(p.PartnerID)
		params.Set("partner_id", val)
	}

	if p.Test {
		params.Set("test", "1")
	}

	if p.Time.After(time.Now()) {
		val := strconv.FormatInt(p.Time.Unix(), 10)
		params.Set("time", val)
	}

	if p.Translit {
		params.Set("translit", "1")
	}

	res, lines, err := c.makeRequest("/sms/send", params)
	if err != nil {
		return Response{}, err
	}

	var ids []string
	re := regexp.MustCompile("^balance=")

	for i := 1; i < len(lines); i++ {
		isBalance := re.MatchString(lines[i])

		if isBalance {
			str := re.ReplaceAllString(lines[i], "")
			balance, err := strconv.ParseFloat(str, 32)
			if err != nil {
				return Response{}, error_internal
			}
			res.Balance = float32(balance)
		} else {
			ids = append(ids, lines[i])
		}
	}

	res.Ids = ids
	return res, nil
}

// SmsStatus will get a status of message
func (c *Client) SmsStatus(id string) (Response, error) {
	params := url.Values{}
	params.Set("id", id)

	res, _, err := c.makeRequest("/sms/status", params)
	if err != nil {
		return Response{}, err
	}

	return res, nil
}

// SmsCost will get a cost of message
func (c *Client) SmsCost(p *Request) (Response, error) {
	var params = url.Values{}
	params.Set("to", p.To)
	params.Set("msg", p.Msg)
	if p.Translit {
		params.Set("translit", "1")
	}

	res, lines, err := c.makeRequest("/sms/cost", params)
	if err != nil {
		return Response{}, err
	}

	cost, err := strconv.ParseFloat(lines[1], 32)
	if err != nil {
		return Response{}, error_internal
	}

	count, err := strconv.Atoi(lines[2])
	if err != nil {
		return Response{}, error_internal
	}

	res.Cost = float32(cost)
	res.Count = count

	return res, nil
}

// MyBalance checks the balance
func (c *Client) MyBalance() (Response, error) {
	res, lines, err := c.makeRequest("/my/balance", url.Values{})
	if err != nil {
		return Response{}, err
	}

	balance, err := strconv.ParseFloat(lines[1], 32)
	if err != nil {
		return Response{}, error_internal
	}

	res.Balance = float32(balance)
	return res, nil
}

// MyLimit checks the limit
// func (c *Client) MyLimit() (Response, error) {
// 	res, lines, err := c.makeRequest("/my/limit", url.Values{})
// 	if err != nil {
// 		return Response{}, err
// 	}

// 	limit, err := strconv.Atoi(lines[1])
// 	if err != nil {
// 		return Response{}, error_internal
// 	}

// 	limitSent, err := strconv.Atoi(lines[2])
// 	if err != nil {
// 		return Response{}, error_internal
// 	}

// 	res.Limit = limit
// 	res.LimitSent = limitSent
// 	return res, nil
// }

// MySenders recieves the list of senders
// func (c *Client) MySenders() (Response, error) {
// 	res, lines, err := c.makeRequest("/my/senders", url.Values{})
// 	if err != nil {
// 		return Response{}, err
// 	}

// 	var senders []string
// 	for i := 1; i < len(lines); i++ {
// 		senders = append(senders, lines[i])
// 	}

// 	res.Senders = senders
// 	return res, nil
// }

// StoplistGet recieves the stoplist
// func (c *Client) StoplistGet() (Response, error) {
// 	res, lines, err := c.makeRequest("/stoplist/get", url.Values{})
// 	if err != nil {
// 		return Response{}, err
// 	}

// 	stoplist := make(map[string]string)
// 	for i := 1; i < len(lines); i++ {
// 		re := regexp.MustCompile(";")
// 		str := re.Split(lines[i], 2)

// 		stoplist[str[0]] = str[1]
// 	}

// 	res.Stoplist = stoplist
// 	return res, nil
// }

// StoplistAdd will add the phone number to stoplist
//
// phone is phone number, text is the additional information.
// func (c *Client) StoplistAdd(phone, text string) (Response, error) {
// 	params := url.Values{}
// 	params.Set("stoplist_phone", phone)
// 	params.Set("stoplist_text", text)

// 	res, _, err := c.makeRequest("/stoplist/add", params)
// 	if err != nil {
// 		return Response{}, err
// 	}

// 	return res, nil
// }

// StoplistDel will delete the phone number from stoplist
//
// phone is phone number
// func (c *Client) StoplistDel(phone string) (Response, error) {
// 	params := url.Values{}
// 	params.Set("stoplist_phone", phone)

// 	res, _, err := c.makeRequest("/stoplist/del", params)
// 	if err != nil {
// 		return Response{}, err
// 	}

// 	return res, nil
// }

// CallbackGet recieves the callbacks from service
// func (c *Client) CallbackGet() (Response, error) {
// 	res, lines, err := c.makeRequest("/callback/get", url.Values{})
// 	if err != nil {
// 		return Response{}, err
// 	}

// 	var callbacks []string
// 	for i := 1; i < len(lines); i++ {
// 		callbacks = append(callbacks, lines[i])
// 	}

// 	res.Callbacks = callbacks
// 	return res, nil
// }

// CallbackAdd will add the callback url to service
//
// cbUrl is your callback url
// func (c *Client) CallbackAdd(cbUrl string) (Response, error) {
// 	params := url.Values{}
// 	params.Set("url", cbUrl)

// 	res, lines, err := c.makeRequest("/callback/add", params)
// 	if err != nil {
// 		return Response{}, err
// 	}

// 	var callbacks []string
// 	for i := 1; i < len(lines); i++ {
// 		callbacks = append(callbacks, lines[i])
// 	}

// 	res.Callbacks = callbacks
// 	return res, nil
// }

// CallbackDel will delete the callback url from service
//
// cbUrl is your callback url
// func (c *Client) CallbackDel(cbUrl string) (Response, error) {
// 	params := url.Values{}
// 	params.Set("url", cbUrl)

// 	res, lines, err := c.makeRequest("/callback/del", params)
// 	if err != nil {
// 		return Response{}, err
// 	}

// 	var callbacks []string
// 	for i := 1; i < len(lines); i++ {
// 		callbacks = append(callbacks, lines[i])
// 	}

// 	res.Callbacks = callbacks
// 	return res, nil
// }
