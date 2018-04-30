package smsru

import (
	"strings"
)

func HairPhone(phone string) string {
	return strings.Replace(phone, "+", "", -1)
}
