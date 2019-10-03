package queue

import (
	"fmt"
	"log"

	"github.com/dstpierre/gosaas/queue/email"
	"github.com/dstpierre/gosaas/internal/config"
)

type SendEmailParameter struct {
	From    string `json:"From"`
	To      string `json:"To"`
	Subject string `json:"Subject"`
	Body    string `json:"Body"`
}

type Email struct {
	Send func(p SendEmailParameter) error
}

type Emailer interface {
	Send(toEmail, toName, fromEmail, fromName, subject, body, replyTo string) error
}

func (e *Email) Run(qt QueueTask) error {
	m, ok := qt.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("this data is not a proper SendEmailParameter")
	}

	var p SendEmailParameter
	if err := fillStruct(&p, m); err != nil {
		return fmt.Errorf("error fill struct: %s", err.Error())
	}

	return e.Send(p)
}

func (e *Email) sendEmailDev(p SendEmailParameter) error {
	fmt.Println("email would have been sent:")
	fmt.Println(p)
	return nil
}

func (e *Email) sendEmailProd(p SendEmailParameter) error {
	var emailer Emailer
	if config.Current.EmailProvider == "amazonses" {
		emailer = &email.AmazonSES{}
	}

	if emailer == nil {
		log.Println("cannot find email provider named: %s", config.Current.EmailProvider)		
		return nil
	}

	return emailer.Send(p.To, p.To, p.From, p.From, p.Subject, p.Body, "")
}
