package push

// Text 文本消息
type Text struct {
	Webhook, AtUserIds, Text string
}

func (t Text) getContent() (url string, body interface{}, header map[string]string) {
	body = H{
		"msgtype": "text",
		"text": H{
			"content": t.Text,
		},
		"at": H{
			"isAtAll":   "False",
			"atUserIds": []string{t.AtUserIds},
		},
	}
	url = t.Webhook
	header = map[string]string{
		"Content-Type": "application/json",
	}
	return
}

// Markdown markdown消息
type Markdown struct {
	Webhook, AtUserIds, Text string
	Title                    string
}

func (t Markdown) getContent() (url string, body interface{}, header map[string]string) {
	body = H{
		"msgtype": "markdown",
		"markdown": H{
			"title": t.Title,
			"text":  t.Text,
		},
		"at": H{
			"isAtAll":   "False",
			"atUserIds": []string{t.AtUserIds},
		},
	}
	header = map[string]string{
		"Content-Type": "application/json",
	}
	url = t.Webhook
	return
}

// ActionCard actionCard消息
type ActionCard struct {
	Webhook, AtUserIds, Text      string
	Title, SingleTitle, SingleURL string
}

func (t ActionCard) getContent() (url string, body interface{}, header map[string]string) {
	body = H{
		"msgtype": "actionCard",
		"actionCard": H{
			"title":       t.Title,
			"text":        t.Text,
			"singleTitle": t.SingleTitle,
			"singleURL":   t.SingleURL,
		},
		"at": H{
			"isAtAll":   "False",
			"atUserIds": []string{t.AtUserIds},
		},
	}
	url = t.Webhook
	header = map[string]string{
		"Content-Type": "application/json",
	}
	return
}
