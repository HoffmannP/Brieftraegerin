package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"time"
)

const NAME = "brieftraegerin"

func main() {
	var configFile, logFile string
	flag.StringVar(&configFile, "config", "config.toml", "Config file")
	flag.StringVar(&configFile, "f", "config.toml", "Config file")
	flag.StringVar(&logFile, "log", "/var/log/"+NAME+".log", "full path of the log file")
	flag.StringVar(&logFile, "l", "/var/log/"+NAME+".log", "full path of the log file")
	flag.Parse()

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("[" + NAME + "] ")
	log.Println("started program")

	/*
		emailSource, err := io.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		mailbkp, err := os.OpenFile("receivedMail.eml", os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}

		pipeRead, pipeWrite := io.Pipe()
		mailbkp, err := os.OpenFile("receivedMail.eml", os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer mailbkp.Close()
		multiWriter := io.MultiWriter(mailbkp, pipeWrite)
		io.Copy(multiWriter, os.Stdin)
		email, err := mail.ReadMessage(pipeRead)
	*/

	email, err := mail.ReadMessage(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("received e-mail from stdin")

	c := readConfig(configFile)
	log.Printf("read config from file %s\n", configFile)

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

	log.Println("sent e-mail to list")
	file.Close()
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
