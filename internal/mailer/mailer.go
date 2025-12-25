package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type InvoiceLinkData struct {
	CustomerName   string
	EstimateNumber int
	SignURL        string
	ExpiresAt      string
}

type Mailer struct {
	dialer    *gomail.Dialer
	fromEmail string
	template  *template.Template
}

func NewMailer() (*Mailer, error) {
	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		return nil, errors.New("SMTP_PORT not set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	tmpl, err := template.ParseFiles("./ui/html/mail/customerInvoiceLink.tmpl")
	if err != nil {
		return nil, err
	}

	return &Mailer{
		dialer: gomail.NewDialer(
			os.Getenv("SMTP_HOST"),
			port,
			os.Getenv("SMTP_USER"),
			os.Getenv("SMTP_PASS"),
		),
		fromEmail: os.Getenv("FROM_EMAIL"),
		template:  tmpl,
	}, nil
}

func (m *Mailer) SendInvoiceLink(to string, data InvoiceLinkData) error {
	log.Println("SendInvoiceLink called for:", to)

	var body bytes.Buffer
	t, err := template.ParseFiles("./ui/html/mail/customerInvoiceLink.tmpl")
	if err != nil {
		return err
	}
	err = t.Execute(&body, data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()

	msg.SetHeader("From", os.Getenv("SMTP_USER"))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", fmt.Sprintf("EzKitchen Invoice Agreement #%d", data.EstimateNumber))
	msg.SetBody("text/html", body.String())

	return m.dialer.DialAndSend(msg)
}
