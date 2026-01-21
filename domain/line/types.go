package line

// ReplyMessage - LINE 답장 메시지 구조체
type ReplyMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ReplyRequest - LINE 답장 요청 구조체
type ReplyRequest struct {
	ReplyToken string         `json:"replyToken"`
	Messages   []ReplyMessage `json:"messages"`
}
