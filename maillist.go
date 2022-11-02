package main

import (
	"fmt"
	"net/mail"
	"strings"
)

type Maillist struct {
	Sender   []string
	Receiver []string
	Server   string
	name     string
	to       string
}

func (l Maillist) CheckValidSender(c Config, from string) (err error) {
	f, err := mail.ParseAddress(from)
	if err != nil {
		return
	}

	for _, validSender := range l.Sender {
		s, _ := mail.ParseAddress(validSender)
		if strings.EqualFold(f.Address, s.Address) {
			return
		}
	}
	err = fmt.Errorf("From %s is no valid sender address", f.Address)
	return
}

func (l Maillist) Recipients() (to []string, err error) {
	for _, r := range l.Receiver {
		t, err := mail.ParseAddress(r)
		if err != nil {
			return to, err
		}
		to = append(to, t.Address)
	}
	return
}
