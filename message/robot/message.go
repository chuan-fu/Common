package robot

// RobotResponse 机器人接口响应
type RobotResponse struct {
	ErrorCode    int64  `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

// Message 消息内容
type Message struct {
	MsgType  string `json:"msgtype"` // 消息类型
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown,omitempty"` // markdown 类型
	Text Text `json:"text,omitempty"` // 文本类型
}

// Text 文本
type Text struct {
	Content             string   `json:"content"`                         // 文本内容，最长不超过2048个字节，必须是utf8编码
	MentionedList       []string `json:"mentioned_list,omitempty"`        // userid的列表，提醒群中的指定成员(@某个成员)，@all表示提醒所有人，如果开发者获取不到userid，可以使用mentioned_mobile_list
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"` // 手机号列表，提醒手机号对应的群成员(@某个成员)，@all表示提醒所有人
}
