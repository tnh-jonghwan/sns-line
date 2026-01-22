package instagram

// WebhookRequest - Instagram 웹훅 요청 구조체
type WebhookRequest struct {
	Object string  `json:"object"` // "instagram"
	Entry  []Entry `json:"entry"`
}

// Entry - 웹훅 엔트리
type Entry struct {
	ID        string      `json:"id"`
	Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

// Messaging - 메시지 정보
type Messaging struct {
	Sender    User     `json:"sender"`
	Recipient User     `json:"recipient"`
	Timestamp int64    `json:"timestamp"`
	Message   *Message `json:"message,omitempty"`
}

// User - 사용자 정보
type User struct {
	ID string `json:"id"`
}

// Message - 메시지 내용
type Message struct {
	Mid  string `json:"mid"`
	Text string `json:"text,omitempty"`
}
