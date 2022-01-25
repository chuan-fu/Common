package email

import "gopkg.in/gomail.v2"

func NewClient(host string, port int, username, password string) *gomail.Dialer {
	return gomail.NewDialer(host, port, username, password)
}

type message struct {
	msg *gomail.Message
}

func NewMessage() *message {
	return &message{msg: gomail.NewMessage()}
}

func (m *message) SetFrom(d string) *message {
	m.msg.SetHeader("From", d)
	return m
}

func (m *message) SetTo(d string) *message {
	m.msg.SetHeader("To", d)
	return m
}

func (m *message) SetSubject(d string) *message {
	m.msg.SetHeader("Subject", d)
	return m
}

func (m *message) SetHtml(d string) *message {
	m.msg.SetBody("text/html", d)
	return m
}

func (m *message) GetMsg() *gomail.Message {
	return m.msg
}

func (m *message) Do(dialer *gomail.Dialer) error {
	return dialer.DialAndSend(m.msg)
}
