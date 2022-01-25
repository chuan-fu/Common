package pushplus

const (
	uri = "https://www.pushplus.plus/send"
)

type (
	TemplateType string
	ChannelType  string
)

const (
	TemplateTypeHtml         TemplateType = "html"         // 支持html文本。为空默认使用html模板
	TemplateTypeTxt          TemplateType = "txt"          // 纯文本内容,不转义html内容,换行使用\n
	TemplateTypeJson         TemplateType = "json"         // 可视化展示json格式内容
	TemplateTypeMarkdown     TemplateType = "markdown"     // 内容基于markdown格式展示
	TemplateTypeCloudMonitor TemplateType = "cloudMonitor" // 阿里云监控报警定制模板
)

const (
	ChannelTypeWechat  ChannelType = "wechat"  // 微信公众号,默认发送渠道
	ChannelTypeWebhook ChannelType = "webhook" // 第三方webhook服务；企业微信机器人、钉钉机器人、飞书机器人
	ChannelTypeCp      ChannelType = "cp"      // 企业微信应用
	ChannelTypeMail    ChannelType = "mail"    // 邮件
	ChannelTypeSms     ChannelType = "sms"     // 短信，未开放使用
)

type Message struct {
	Token       string       `json:"token"`                  // 用户令牌
	Title       string       `json:"title,omitempty"`        // 消息标题
	Content     string       `json:"content"`                // 具体消息内容，根据不同template支持不同格式
	Template    TemplateType `json:"template,omitempty"`     // 发送消息模板 默认html
	Channel     ChannelType  `json:"channel,omitempty"`      // 发送渠道 默认wechat
	Webhook     string       `json:"webhook,omitempty"`      // webhook编码，仅在channel使用webhook渠道和CP渠道时需要填写
	CallbackUrl string       `json:"callback_url,omitempty"` // 回调地址，异步回调发送结果 非必填
	Timestamp   string       `json:"timestamp,omitempty"`    // 时间戳，毫秒。如小于当前时间，消息将无法发送 非必填
	Topic       string       `json:"topic,omitempty"`        // 群组编码
}

type Result struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
