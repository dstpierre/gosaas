package queue

import (
	"fmt"
	"log"
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
	log.Fatal("not implemented")
	return nil
}
