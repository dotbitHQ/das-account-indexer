package toolib

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

type EmailHelper struct {
	Host           string
	Port           int
	From           string
	SenderAddress  string
	SenderPassword string
}

func (e *EmailHelper) SendEmail(subject, body string, toList ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From+"<"+e.SenderAddress+">")
	m.SetHeader("To", toList...)
	m.SetHeader("Content-Type", "text/html; charset=UTF-8")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(e.Host, e.Port, e.SenderAddress, e.SenderPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d.DialAndSend(m)
}
