package test

import (
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type ImapMail struct {
	EmailAddress string
	Password     string
	Client       *client.Client
}

func NewImapMail(email string, psw string) *ImapMail {
	return &ImapMail{
		EmailAddress: email,
		Password:     psw,
	}
}

func (m *ImapMail) Login() error {
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	//defer c.Logout()

	// Login
	if err := c.Login(m.EmailAddress, m.Password); err != nil {
		log.Fatal(err)
		return err
	}
	m.Client = c
	return nil
}
func (m *ImapMail) GetNewMail() (string, error) {
	c := m.Client
	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// Get the last message
	if mbox.Messages == 0 {
		log.Println("No messages in mailbox")
		return "", nil
	}

	// seqset := new(imap.SeqSet)
	// seqset.AddRange(mbox.Messages, mbox.Messages)

	// messages := make(chan *imap.Message, 1)
	// done := make(chan error, 1)
	// go func() {
	// 	done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	// }()

	// msg := <-messages
	// if msg == nil {
	// 	log.Println("Server didn't return message")
	// 	return "", nil
	// }

	// log.Println("Got message:", msg.Envelope.Subject)
	// log.Println("Date:", msg.Envelope.Date.Format(time.RFC3339))
	// return msg.Envelope.Subject, nil
	// Get the last 10 messages
	lastNum := uint32(10)

	seqset := new(imap.SeqSet)
	seqset.AddRange(mbox.Messages-lastNum+1, mbox.Messages)

	messages := make(chan *imap.Message, lastNum)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	for i := 0; i < int(lastNum); i++ {
		msg := <-messages
		if msg == nil {
			log.Println("Server didn't return message")
			break
		}

		log.Println("Got message:", msg.Envelope.From[0], msg.Envelope.Subject)
		//log.Println("Date:", msg.Envelope.Date.Format(time.RFC3339))
	}
	return " oooe", nil
}
