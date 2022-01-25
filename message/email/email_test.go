package email

import (
	"context"
	"testing"

	"gopkg.in/gomail.v2"
)

func TestSendEmail(t *testing.T) {
	host, port, username, password := "smtp.qq.com", 587, "xxx", "xxx"
	type args struct {
		ctx     context.Context
		dialer  *gomail.Dialer
		from    string
		to      string
		subject string
		body    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "text推送测试",
			args: args{
				ctx:     context.Background(),
				dialer:  gomail.NewDialer(host, port, username, password),
				from:    "xxx",
				to:      "xxx",
				subject: "主题1",
				body:    "哈哈哈嘿嘿嘿",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewMessage().SetFrom(tt.args.from).SetTo(tt.args.to).SetSubject(tt.args.subject).SetHtml(tt.args.body).Do(tt.args.dialer); (err != nil) != tt.wantErr {
				t.Errorf("SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
