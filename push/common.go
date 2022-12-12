package push

// Common 通用消息接口
type Common interface {
	getContent() (string, interface{}, map[string]string)
}

// Ding 钉钉消息
type Ding struct {
	ConversationId string `json:"conversationId"`
	SenderStaffId  string `json:"senderStaffId"`
	AtUsers        []struct {
		DingtalkId string `json:"dingtalkId"`
	} `json:"atUsers"`
	ChatbotUserId             string `json:"chatbotUserId"`
	MsgId                     string `json:"msgId"`
	SenderNick                string `json:"senderNick"`
	IsAdmin                   bool   `json:"isAdmin"`
	SessionWebhookExpiredTime int64  `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64  `json:"createAt"`
	ConversationType          string `json:"conversationType"`
	SenderId                  string `json:"senderId"`
	ConversationTitle         string `json:"conversationTitle"`
	IsInAtList                bool   `json:"isInAtList"`
	SessionWebhook            string `json:"sessionWebhook"`
	Text                      struct {
		Content string `json:"content"`
	} `json:"text"`
	RobotCode string `json:"robotCode"`
	Msgtype   string `json:"msgtype"`
}

// Result 发送消息请求结果
type Result struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type H map[string]interface{}
