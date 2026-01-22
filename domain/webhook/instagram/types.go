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
	Sender    User      `json:"sender"`
	Recipient User      `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   *Message  `json:"message,omitempty"`
	Postback  *Postback `json:"postback,omitempty"`
	Read      *Read     `json:"read,omitempty"`
	Delivery  *Delivery `json:"delivery,omitempty"`
}

// User - 사용자 정보
type User struct {
	ID string `json:"id"`
}

// Message - 메시지 내용
type Message struct {
	Mid     string `json:"mid"`
	Text    string `json:"text,omitempty"`
	IsEcho  bool   `json:"is_echo,omitempty"`
	ReplyTo *struct {
		Mid string `json:"mid"`
	} `json:"reply_to,omitempty"`
}

// Postback - 버튼 클릭 등의 이벤트
type Postback struct {
	Title   string `json:"title"`
	Payload string `json:"payload"`
}

// Read - 읽음 확인
type Read struct {
	Watermark int64 `json:"watermark"`
}

// Delivery - 전달 확인
type Delivery struct {
	Mids      []string `json:"mids"`
	Watermark int64    `json:"watermark"`
}
