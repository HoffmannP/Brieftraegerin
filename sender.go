package main

import (
	"fmt"
	"net/smtp"
)

type SmtpServer struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

func (s SmtpServer) send(to []string, body []byte) error {
	if len(s.From) == 0 {
		s.From = s.User
	}

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.Host, s.Port),
		smtp.PlainAuth("", s.User, s.Password, s.Host),
		s.From, to, body)
}
