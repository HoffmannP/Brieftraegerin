package main

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Domain  string
	Address string
	Test    bool
	Smtp    map[string]SmtpServer
	List    map[string]Maillist
}

func readConfig(filepath string) (config Config) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	decoder := toml.NewDecoder(f)
	_, err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func (c Config) selectList(e mail.Message) (l Maillist, err error) {
	to, err := mail.ParseAddress(e.Header.Get("To"))
	if err != nil {
		return
	}
	pattern, err := regexp.Compile("^" + strings.Replace(c.Address, "%", "([\\w-]{3,})", 1) + "@" + c.Domain + "$")
	if err != nil {
		return
	}
	name := pattern.FindStringSubmatch(to.Address)[1]
	if len(name) < 3 {
		err = fmt.Errorf("Target %s does not match format %s", c.Address+"@"+c.Domain, pattern)
		return
	}
	l, ok := c.List[name]
	if !ok {
		err = fmt.Errorf("List %s was not found", name)
		return
	}

	if l.Server == "" {
		l.Server = "default"
	}
	l.name = name
	l.to = newTo(c, l)

	err = l.CheckValidSender(c, e.Header.Get("From"))
	return
}

func (c Config) selectServer(l Maillist) (s SmtpServer, err error) {
	s, ok := c.Smtp[l.Server]
	if !ok {
		for _, s_ := range c.Smtp {
			s = s_
			ok = true
			break
		}
	}
	if !ok {
		err = errors.New("No SMTP server defined")
	}
	return
}

func newTo(c Config, l Maillist) string {
	pattern := c.Address + "@" + c.Domain
	format := strings.Replace(pattern, "%", "%s", 1)
	return fmt.Sprintf(format, l.name)
}
