package smsru_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oanogin/go-smsru"
)

func getPhone() string {
	return os.Getenv("phone")
}

func getClient(t *testing.T) *smsru.Client {
	apiId := os.Getenv("api_id")
	return smsru.NewClient(apiId, true, true, false)
}

/* Test Sms
---------------------------------------------*/
func TestSmsSend(t *testing.T) {
	c := getClient(t)
	phone := getPhone()

	ssr, err := c.SendSms(phone, "Sample")

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 100, ssr.StatusCode)
	assert.Equal(t, 100, ssr.Sms[phone].StatusCode)

}
