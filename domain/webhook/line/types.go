package line

type WebhookRequest struct {
	Destination string  `json:"destination"`
	Events      []Event `json:"events"`
}

type Event struct {
	Type       string   `json:"type"`
	Timestamp  int64    `json:"timestamp"`
	Source     Source   `json:"source"`
	ReplyToken string   `json:"replyToken,omitempty"`
	Message    *Message `json:"message,omitempty"`
}

type Source struct {
	Type   string `json:"type"`
	UserID string `json:"userId,omitempty"`
}

type Message struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}
