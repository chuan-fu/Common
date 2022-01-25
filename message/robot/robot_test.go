package robot

import (
	"context"
	"testing"
)

const key = ""

func TestRobot_SendMarkdown(t *testing.T) {
	type args struct {
		ctx      context.Context
		key      string
		markdown string
		userId   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "markdown推送测试",
			args: args{
				ctx: context.Background(),
				key: key,
				markdown: `实时新增用户反馈<font color=\"warning\">132例</font>，请相关同事注意。
										>类型:<font color=\"comment\">用户反馈</font>
         								>普通用户反馈:<font color=\"comment\">117例</font>
                                        >VIP用户反馈:<font color=\"comment\">15例</font>`,
				userId: []string{"<@xxx>"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rbt := NewRobot()
			if err := rbt.SendMarkdown(tt.args.ctx, tt.args.key, tt.args.markdown, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("SendMarkdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRobot_SendText(t *testing.T) {
	type args struct {
		ctx                 context.Context
		key                 string
		text                string
		mentionedList       []string
		mentionedMobileList []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "text推送测试",
			args: args{
				ctx:           context.Background(),
				key:           key,
				text:          "大家好，这是我的一个企微群机器人测试",
				mentionedList: []string{"@all"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rbt := NewRobot()
			if err := rbt.SendText(tt.args.ctx, tt.args.key, tt.args.text, tt.args.mentionedList, tt.args.mentionedMobileList); (err != nil) != tt.wantErr {
				t.Errorf("SendText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
