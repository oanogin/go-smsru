package smsru

import (
	"strings"
)

func hairPhone(phone string) string {
	return strings.Replace(phone, "+", "", -1)
}
