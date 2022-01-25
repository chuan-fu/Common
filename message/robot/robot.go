package robot

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type robot struct{}

func NewRobot() *robot {
	return &robot{}
}

// SendText 发送文本消息
func (rbt *robot) SendText(ctx context.Context, key, text string, mentionedList, mentionedMobileList []string) (err error) {
	if len(strings.TrimSpace(key)) == 0 {
		return errors.New("无效的企微机器人key")
	}

	data := Message{
		MsgType: "text",
		Text: Text{
			Content:             text,
			MentionedList:       mentionedList,
			MentionedMobileList: mentionedMobileList,
		},
	}
	if err = post(ctx, key, data); err != nil {
		return err
	}
	return nil
}

// SendMarkdown 发送markdown消息
func (rbt *robot) SendMarkdown(ctx context.Context, key, markdown string, userIds []string) (err error) {
	if len(strings.TrimSpace(key)) == 0 {
		return errors.New("无效的企微机器人key")
	}

	data := Message{
		MsgType: "markdown",
		Markdown: struct {
			Content string `json:"content"`
		}{Content: func() string {
			if len(userIds) == 0 {
				return markdown
			}

			return fmt.Sprintf("%s\n\n%s", markdown, strings.Join(userIds, " "))
		}()},
	}

	if err = post(ctx, key, data); err != nil {
		return err
	}
	return nil
}
