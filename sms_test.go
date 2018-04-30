package smsru_test

import (
	"os"
	"testing"

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
	phone := smsru.HairPhone(getPhone())

	_, err := c.SendSms(phone, "Sample")

	if err != nil {
		t.Fail()
	}

}
