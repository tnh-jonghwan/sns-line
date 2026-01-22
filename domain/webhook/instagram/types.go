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
	Messaging []Messaging `json:"messaging,omitempty"` // DM 메시지용
	Changes   []Change    `json:"changes,omitempty"`   // 댓글, 좋아요 등
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

// Change - Instagram changes 이벤트 (댓글, 좋아요 등)
type Change struct {
	Field string      `json:"field"` // "comments", "feed", "mentions" 등
	Value ChangeValue `json:"value"`
}

// ChangeValue - Change 안의 실제 데이터
type ChangeValue struct {
	From struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"from,omitempty"`
	Media struct {
		ID               string `json:"id"`
		MediaProductType string `json:"media_product_type"`
	} `json:"media,omitempty"`
	ID       string `json:"id"`
	ParentID string `json:"parent_id,omitempty"`
	Text     string `json:"text,omitempty"`
}
