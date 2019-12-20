package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

var auth smtp.Auth

//Request struct
type Email struct {
	refEmail  string
	password  string
	smtpEmail string
}

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func (email *Email) initEmail() {
	// SMTP "smtp.gmail.com"
	auth = smtp.PlainAuth("", email.refEmail, email.password, email.smtpEmail)
}

func NewEmail(refEmail string, password string, smtpEmail string) *Email {
	return &Email{
		refEmail:  refEmail,
		password:  password,
		smtpEmail: smtpEmail,
	}
}

func NewEmailRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (email *Email) SendEmail(req *Request, tmplPath string, data interface{}) (bool, error) {
	fmt.Println("Initial email...")
	email.initEmail()

	fmt.Println("Start parsing html template...")
	err := req.parseTemplate(tmplPath, data)
	if err != nil {
		return false, fmt.Errorf("Error in parsing html template: %s", err.Error())
	}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + req.subject + "!\n"
	msg := []byte(subject + mime + "\n" + req.body)
	addr := "smtp.gmail.com:587"

	fmt.Printf("\n********************\nStart sending email from %s to %s...\n********************\n", email.refEmail, req.to)
	errEmail := smtp.SendMail(addr, auth, email.refEmail, req.to, msg)
	if errEmail != nil {
		return false, fmt.Errorf("Error sending email: %s", errEmail.Error())
	}
	fmt.Printf("\n********************\nEmail sent from %s to %s...\n********************\n", email.refEmail, req.to)
	return true, nil
}

func (r *Request) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	fmt.Printf("HTML parsed %s   __---__   data: %v", r.body, data)
	return nil
}
