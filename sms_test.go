package smsru_test

import (
	"fmt"
	"github.com/oanogin/go-smsru"
	"os"
	"testing"
)

func getPhone() string {
	return os.Getenv("phone")
}

func getClient(t *testing.T) *smsru.Client {
	apiId := os.Getenv("api_id")
	return smsru.NewClient(apiId, false, true, false)
}

/* Test Sms
---------------------------------------------*/
func TestSmsSend(t *testing.T) {
	c := getClient(t)
	phone := smsru.HairPhone(getPhone())

	_, err := c.SendSms(phone, fmt.Sprintf("Привет, это прикол, %s", "1234"))

	if err != nil {
		t.Fail()
	}

}
