package main

import (
	"fmt"
	"log"
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
	log.Printf("sending to %d list entries", len(to))

	to = []string{"p2lebe@uni-jena.de"}
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.Host, s.Port),
		smtp.PlainAuth("", s.User, s.Password, s.Host),
		s.From, to, body)
}
