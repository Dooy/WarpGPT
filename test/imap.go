package test

import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
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
		return "", errors.New("No messages in mailbox")
	}

	lastNum := uint32(1)

	seqset := new(imap.SeqSet)
	seqset.AddRange(mbox.Messages-lastNum+1, mbox.Messages)

	messages := make(chan *imap.Message, lastNum)
	done := make(chan error, 1)
	section := &imap.BodySectionName{}
	//items := []imap.FetchItem{section.FetchItem()}
	go func() {
		//done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages)
	}()

	for i := 0; i < int(lastNum); i++ {
		msg := <-messages
		if msg == nil {
			log.Println("Server didn't return message")
			break
		}

		log.Println("Got message:", msg.Envelope.To[0].Address(), msg.Envelope.Subject)

		log.Println("Date:", msg.Envelope.Date.Format(time.RFC3339))

		reader := msg.GetBody(section)
		if reader == nil {
			log.Fatal("未能获取邮件正文")
			return "", nil
		}
		mr, err := mail.CreateReader(reader)
		if err != nil {
			log.Fatal(err)
		}
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				b, _ := io.ReadAll(p.Body)
				log.Printf("Got text: %v\n", string(b))
			case *mail.AttachmentHeader:
				filename, _ := h.Filename()
				log.Printf("Got attachment: %v\n", filename)
			}
		}
	}
	return "", nil
}

func readMail(mr *imap.Reader) {

}
