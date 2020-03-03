package sendgrid

import (
	"encoding/base64"
	"errors"
	"github.com/codingbeard/cbmail"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"strconv"
	"time"
)

type Provider struct {
	config       cbmail.Config
	logger       cbmail.Logger
	errorHandler cbmail.ErrorHandler
}

type Email struct {
	provider        *Provider
	mail            *mail.SGMailV3
	personalization *mail.Personalization
}

func New(dependencies cbmail.Dependencies) (*Provider, error) {
	_, e := dependencies.Config.GetRequiredString("mail.sendgrid.key")
	if e != nil {
		return nil, e
	}

	return &Provider{
		config:       dependencies.Config,
		logger:       dependencies.Logger,
		errorHandler: dependencies.ErrorHandler,
	}, nil
}

func (p *Provider) New() cbmail.Email {
	email := mail.NewV3Mail()
	email.Headers = make(map[string]string)

	return &Email{
		provider: p,
		mail:     email,
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
	if m.personalization == nil {
		return errors.New("no recipients specified")
	}
	if m.mail.From == nil {
		return errors.New("no sender specified")
	}

	m.mail.AddPersonalizations(m.personalization)

	key, e := m.provider.config.GetRequiredString("mail.sendgrid.key")

	if e != nil {
		return e
	}

	trackingSettings := mail.NewTrackingSettings()
	clickTrackingSettings := mail.NewClickTrackingSetting()
	clickTrackingSettings.SetEnable(false)
	clickTrackingSettings.SetEnableText(false)
	trackingSettings.SetClickTracking(clickTrackingSettings)
	openTrackingSetting := mail.NewOpenTrackingSetting()
	openTrackingSetting.SetEnable(false)
	trackingSettings.SetOpenTracking(openTrackingSetting)
	subscriptionTrackingSetting := mail.NewSubscriptionTrackingSetting()
	subscriptionTrackingSetting.SetEnable(false)
	trackingSettings.SetSubscriptionTracking(subscriptionTrackingSetting)
	googleAnalyticsSetting := mail.NewGaSetting()
	googleAnalyticsSetting.SetEnable(false)
	trackingSettings.SetGoogleAnalytics(googleAnalyticsSetting)
	m.mail.SetTrackingSettings(trackingSettings)

	m.mail.Headers["X-Entity-Ref-ID"] = strconv.FormatInt(time.Now().Unix(), 10)

	request := sendgrid.GetRequest(key, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m.mail)
	_, e = sendgrid.API(request)
	if e != nil {
		return e
	}

	return nil
}
