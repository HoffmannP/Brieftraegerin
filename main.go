package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"time"
)

func main() {
	email, err := mail.ReadMessage(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Read E-Mail from StdIn")

	c := readConfig()
	log.Println("Read Config from File")

	list, err := c.selectList(*email)
	if err != nil {
		log.Fatal(err)
	}

	server, err := c.selectServer(list)
	if err != nil {
		log.Fatal(err)
	}

	err = sendMail(server, list, *email)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sent E-Mail to STMP")
}

func sendMail(s SmtpServer, l Maillist, e mail.Message) (err error) {
	e.Header["Date"] = []string{time.Now().Format(time.RFC1123Z)}
	e.Header["To"] = []string{l.to}
	e.Header["Subject"] = []string{"[" + l.name + "] " + e.Header.Get("Subject")}

	to, err := l.Recipients()
	if err != nil {
		log.Fatal(err)
	}

	return s.send(to, mailToString(e))
}

func mailToString(e mail.Message) []byte {
	msg := make([]byte, 0)
	for key, values := range e.Header {
		for _, value := range values {
			msg = append(msg, fmt.Sprintf("%s: %s\r\n", key, value)...)
		}
	}
	message := bytes.NewBuffer(msg)
	_, err := io.Copy(message, e.Body)
	if err != nil {
		log.Fatal(err)
	}
	return message.Bytes()
}
