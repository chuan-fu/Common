package pushplus

import (
	"context"
	"testing"
)

const token = ""

func TestSendHtml(t *testing.T) {
	type args struct {
		ctx                          context.Context
		token, title, content, topic string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "html推送测试",
			args: args{
				ctx:     context.Background(),
				token:   token,
				title:   "测试消息",
				content: "测试内容",
				topic:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPushPlus()
			if err := p.SendMsg(tt.args.ctx, &Message{
				Token:       tt.args.token,
				Title:       tt.args.title,
				Content:     tt.args.content,
				Template:    "",
				Channel:     "mail",
				Webhook:     "",
				CallbackUrl: "",
				Timestamp:   "",
				Topic:       tt.args.topic,
			}); (err != nil) != tt.wantErr {
				t.Errorf("SendHtml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
