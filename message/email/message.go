package email

import (
	"gopkg.in/gomail.v2"
)

type message struct {
	msg *gomail.Message
}

type Option func(m *message)

// 发件人
func WithFrom(from string) Option {
	return func(m *message) {
		m.msg.SetHeader("From", from)
	}
}

// 收件人
func WithTo(to string) Option {
	return func(m *message) {
		m.msg.SetHeader("To", to)
	}
}

// 邮件标题
func WithSubject(subject string) Option {
	return func(m *message) {
		m.msg.SetHeader("Subject", subject)
	}
}

// 邮件内容
func WithBody(body string) Option {
	return func(m *message) {
		m.msg.SetBody("text/html", body)
	}
}

// 附件
func WithAttach(filename string, settings ...gomail.FileSetting) Option {
	return func(m *message) {
		m.msg.Attach(filename, settings...)
	}
}

// 不是脱裤子放屁
// 是封了一层通用消息结构，包括from,to之类的
// 可以不用
func NewMessage(opts ...Option) *gomail.Message {
	m := &message{msg: gomail.NewMessage()}
	for _, opt := range opts {
		opt(m)
	}
	return m.msg
}

func NewMessageWithoutDefault(msg *gomail.Message, opts ...Option) *gomail.Message {
	if msg == nil {
		msg = gomail.NewMessage()
	}
	m := &message{msg: msg}
	for _, opt := range opts {
		opt(m)
	}
	return m.msg
}
