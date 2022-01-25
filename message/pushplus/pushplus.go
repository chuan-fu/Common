package pushplus

import (
	"context"
)

type PushPlus struct{}

func NewPushPlus() *PushPlus {
	return &PushPlus{}
}

func (p *PushPlus) SendWechatTxt(ctx context.Context, token, title, content string, topic ...string) error {
	if err := post(ctx, &Message{
		Token:    token,
		Title:    title,
		Content:  content,
		Template: TemplateTypeTxt,
		Channel:  ChannelTypeWechat,
		Topic: func() string {
			if len(topic) > 0 {
				return topic[0]
			}
			return ""
		}(),
	}); err != nil {
		return err
	}
	return nil
}

func (p *PushPlus) SendMsg(ctx context.Context, data *Message) error {
	if err := post(ctx, data); err != nil {
		return err
	}
	return nil
}
