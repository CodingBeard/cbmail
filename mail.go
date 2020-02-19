package cbmail

type Provider interface {
	New() Email
}

type Contact struct {
	name string
	email string
}

type Email interface {
	SetFrom(contact *Contact)
	AddTo(contact *Contact)
	AddCC(contact *Contact)
	AddBCC(contact *Contact)
	SetReplyTo(contact *Contact)
	SetSubject(subject string)
	SetHeader(key, value string)
	SetTextBody(body string)
	SetHtmlBody(body string)
	AddAttachment(filename, contentType string, content []byte)
	Send() error
}

func NewContact(name, email string) *Contact {
	return &Contact{
		name:  name,
		email: email,
	}
}

func (c *Contact) GetName() string {
	return c.name
}

func (c *Contact) GetEmail() string {
	return c.email
}