package webhook

// WebhookRequest - LINE에서 보내는 Webhook 요청 구조체
type WebhookRequest struct {
	Destination string  `json:"destination"`
	Events      []Event `json:"events"`
}

// Event - LINE 이벤트 구조체
type Event struct {
	Type       string   `json:"type"`
	Timestamp  int64    `json:"timestamp"`
	Source     Source   `json:"source"`
	ReplyToken string   `json:"replyToken,omitempty"`
	Message    *Message `json:"message,omitempty"`
}

// Source - 이벤트 발생 소스 정보
type Source struct {
	Type   string `json:"type"`
	UserID string `json:"userId,omitempty"`
}

// Message - 메시지 정보
type Message struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ReplyRequest - LINE Reply API 요청 구조체
type ReplyRequest struct {
	ReplyToken string         `json:"replyToken"`
	Messages   []ReplyMessage `json:"messages"`
}

// ReplyMessage - 답장 메시지 구조체
type ReplyMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
