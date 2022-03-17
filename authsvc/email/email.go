package email

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/parthoshuvo/authsvc/cfg"
	log "github.com/parthoshuvo/authsvc/log4u"
	"github.com/parthoshuvo/authsvc/table/user"
)

type Mail struct {
	Sender    user.Email
	Recipient user.Email
	Subject   string
	Body      string
}

func (mail *Mail) buildMessage() (string, error) {
	buf := bytes.NewBuffer([]byte{})
	err := mail.template().Execute(buf, mail)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (mail *Mail) template() *template.Template {
	tmpl := `
	From: {{.Sender}}
	To: {{.Recipient}}
	Subject: {{.Subject}}
	{{.Body}}
	`
	return template.Must(template.New("").Parse(tmpl))
}

type EmailClient struct {
	smtpDef *cfg.SmtpServerDef
	client  *smtp.Client
}

func NewEmailClient(smtpDef *cfg.SmtpServerDef) *EmailClient {
	return &EmailClient{smtpDef, openConnection(smtpDef)}
}

func openConnection(def *cfg.SmtpServerDef) *smtp.Client {
	client, err := smtp.Dial(fmt.Sprintf("%s:%d", def.Host, def.Port))
	if err != nil {
		log.Fatalf("failed to open smtp server %s:%d: [%v]", def.Host, def.Port, err)
	}
	if err := sayHello(client, def.From); err != nil {
		log.Fatalf("failed to send mail from %s: [%v]", def.From, err)
	}
	log.Info("Successfully connected to smtp server!!")
	return client
}

func (c *EmailClient) NewMail(recipient user.Email, subject, message string) *Mail {
	return &Mail{
		Sender:    c.smtpDef.From,
		Recipient: recipient,
		Subject:   subject,
		Body:      message,
	}
}

func (c *EmailClient) SendEmail(mail *Mail) error {
	return sendMail(c.client, mail)
}

func sayHello(c *smtp.Client, from user.Email) error {
	return sendMail(c, &Mail{
		Sender:    from,
		Recipient: from,
		Subject:   "Hello from service",
		Body:      "Hi, I'm online",
	})
}

func sendMail(c *smtp.Client, mail *Mail) error {
	if mail.Sender.IsEmpty() {
		return errors.New("can't send mail coz sender is empty")
	}
	c.Mail(mail.Sender.String())

	if mail.Recipient.IsEmpty() {
		return errors.New("can't send mail coz recipient is empty")
	}
	c.Rcpt(mail.Recipient.String())

	if mail.Body == "" {
		return errors.New("can't send mail coz message body is empy")
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to send mail from %s: [%v]", mail.Sender, err)
	}
	defer w.Close()
	msg, err := mail.buildMessage()
	if err != nil {
		return fmt.Errorf("failed to build message: [%v]", err)
	}
	_, err = w.Write([]byte(msg))
	return err
}

func (svc *EmailClient) Close() {
	svc.client.Close()
}
