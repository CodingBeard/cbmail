package sendgrid

import (
	"encoding/base64"
	"github.com/codingbeard/cbmail"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Provider struct {
	config       cbmail.Config
	logger       cbmail.Logger
	errorHandler cbmail.ErrorHandler
}

type Email struct {
	provider        Provider
	mail            *mail.SGMailV3
	personalization *mail.Personalization
}

func New(dependencies cbmail.Dependencies) *Provider {
	_, e := dependencies.Config.GetRequiredString("mail:sendgrid:key")
	if e != nil {
		return e
	}

	return &Provider{
		config:       dependencies.Config,
		logger:       dependencies.Logger,
		errorHandler: dependencies.ErrorHandler,
	}
}

func (p *Provider) New() cbmail.Email {
	return &Email{
		provider: *p,
		mail:     mail.NewV3Mail(),
	}
}

func (m *Email) SetFrom(contact *cbmail.Contact) {
	m.mail.SetFrom(mail.NewEmail(contact.GetName(), contact.GetEmail()))
}

func (m *Email) AddTo(contact *cbmail.Contact) {
	if m.personalization == nil {
		m.personalization = mail.NewPersonalization()
	}
	m.personalization.AddTos(mail.NewEmail(contact.GetName(), contact.GetEmail()))
}

func (m *Email) AddCC(contact *cbmail.Contact) {
	if m.personalization == nil {
		m.personalization = mail.NewPersonalization()
	}
	m.personalization.AddCCs(mail.NewEmail(contact.GetName(), contact.GetEmail()))
}

func (m *Email) AddBCC(contact *cbmail.Contact) {
	if m.personalization == nil {
		m.personalization = mail.NewPersonalization()
	}
	m.personalization.AddBCCs(mail.NewEmail(contact.GetName(), contact.GetEmail()))
}

func (m *Email) SetReplyTo(contact *cbmail.Contact) {
	m.mail.SetReplyTo(mail.NewEmail(contact.GetName(), contact.GetEmail()))
}

func (m *Email) SetSubject(subject string) {
	m.mail.Subject = subject
}

func (m *Email) SetHeader(key, value string) {
	m.mail.SetHeader(key, value)
}

func (m *Email) SetTextBody(body string) {
	m.mail.AddContent(mail.NewContent("text/plain", body))
}

func (m *Email) SetHtmlBody(body string) {
	m.mail.AddContent(mail.NewContent("text/html", body))
}

func (m *Email) AddAttachment(filename, contentType string, content []byte) {
	a := mail.NewAttachment()
	a.SetFilename(filename)
	a.SetType(contentType)
	a.SetContent(base64.StdEncoding.EncodeToString(content))
	m.mail.AddAttachment(a)
}

func (m *Email) Send() error {
	m.mail.AddPersonalizations(m.personalization)

	key, e := m.provider.config.GetRequiredString("mail:sendgrid:key")

	if e != nil {
		return e
	}

	request := sendgrid.GetRequest(key, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m.mail)
	_, e = sendgrid.API(request)
	if e != nil {
		return e
	}

	return nil
}
