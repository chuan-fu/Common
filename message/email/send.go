package email

import (
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

type sendService struct {
	dialer *gomail.Dialer
}

type SendOption func(s *sendService)

func WithClient(client *gomail.Dialer) SendOption {
	return func(s *sendService) {
		s.dialer = client
	}
}

func NewSendService(opts ...SendOption) *sendService {
	s := &sendService{}
	for _, opt := range opts {
		opt(s)
	}
	if s.dialer == nil {
		s.dialer = globalClient
	}
	return s
}

func (s *sendService) SendMessage(msgs ...*gomail.Message) error {
	if s.dialer == nil {
		return errors.New("client为空，无法发送")
	}
	return s.dialer.DialAndSend(msgs...)
}
